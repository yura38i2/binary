package binary

import (
	"encoding/binary"
	"errors"
	"fmt"
	"math"
	"reflect"
)

// Package binary предоставляет функционал для сериализации и десериализации Go-объектов в бинарный формат.
// Используется little-endian порядок байт для всех числовых типов.
// Поддерживаемые типы: bool, int/uint (все размеры), float32/64, string, slice, array, struct.

const (

	// Важно: Намерено указано 1200, потому что для данной задачи
	// объем кодируемых данных должен помещаться в 1 KCP пакет (DataGram pack)
	//
	// Почему не 1400: с заголовками IP/UDP датаграмма в 1404 байта — это ~1432 байта на
	// проводе. В мобильных сетях (особенно IPv6, минимум 1280) и за VPN такие пакеты
	// фрагментируются или пропадают, причём «полноразмерные» теряются, а короткие проходят
	// — тяжёлый для диагностики симптом. 1204 + 28 = 1232 проходит везде.
	// 4 байта - это CRC проверка. Остается 1200.
	// 1200 - с запасом хватает для данного проекта.
	// startBufSize - стартовая велична буфера кодирования
	startBufSize = 1200

	// maxStringLength - максимальная разрешенная длина строки в сериализованных данных.
	// Ограничение защищает от exhaustion памяти при обработке очень длинных строк.
	maxStringLength = uint32(1 << 20) // 1M элементов, не MB данных

	// maxSliceLength - максимальное количество элементов в слайсе/массиве для сериализации.
	// Ограничение предотвращает создание чрезмерно больших объектов в памяти.
	maxSliceLength = uint32(1 << 20) // 1M элементов, не MB данных

	// nilSliceSentinel - специальное значение (все биты установлены) для кодирования nil слайсов.
	// Используется для различения nil слайса и пустого слайса в бинарном представлении.
	nilSliceSentinel = uint32(0xFFFFFFFF) // Все биты установлены, явно отличная метка-сентинел
)

var (
	// ErrUnsupportedType возвращается при попытке сериализации или десериализации неподдерживаемого типа.
	ErrUnsupportedType = errors.New("unsupported type")

	// ErrInvalidBoolValue возвращается при декодировании значения bool, которое не равно 0 или 1.
	ErrInvalidBoolValue = errors.New("invalid boolean value: must be 0 or 1")

	// ErrTruncatedData возвращается когда входные данные недостаточны для полной десериализации.
	ErrTruncatedData = errors.New("truncated data: insufficient bytes available")

	// ErrExtraBytes возвращается когда после успешной десериализации остаются лишние байты в данных.
	ErrExtraBytes = errors.New("extra bytes remaining after decoding")

	// ErrInterfaceTypeUnknown возвращается когда тип интерфейса не может быть определен из сериализованных данных.
	// Пользователю следует предоставить конкретный тип назначения вместо interface{}.
	ErrInterfaceTypeUnknown = errors.New("interface{} contains unknown marshaled type; provide a concrete destination type instead")
)

// Marshal сериализует Go-значение в его детерминированное бинарное представление.
// Формат: каждый тип кодируется с использованием little-endian порядка байт.
// Структуры сериализуются по порядку объявления полей (включая неэкспортируемые поля).
func Marshal(v any) ([]byte, error) {
	if v == nil {
		return nil, fmt.Errorf("cannot marshal nil value")
	}

	vv := reflect.ValueOf(v)

	// Handle interface containing a concrete value (e.g., any(42))
	if vv.Kind() == reflect.Interface && vv.Elem().IsValid() {
		return Marshal(vv.Elem().Interface())
	}

	// Allow nil slice pointers to be marshaled as sentinel values
	if vv.Kind() == reflect.Ptr && vv.IsNil() {
		// Only reject nil pointers that aren't slices
		if vv.Type().Elem().Kind() != reflect.Slice {
			return nil, fmt.Errorf("cannot marshal nil value")
		}
		// FIX: Initialize nil slice pointer to empty slice before marshaling
		elemType := vv.Type().Elem()
		vv = reflect.MakeSlice(elemType, 0, 0)
	}

	encoder := Encoder{buf: make([]byte, 0, startBufSize)}
	err := encoder.marshalValue(vv)
	if err != nil {
		return nil, err
	}

	return encoder.buf, nil
}

// Encoder строит бинарное представление значений.
type Encoder struct {
	buf []byte
}

func newEncoder() *Encoder {
	return &Encoder{buf: make([]byte, 0, 256)}
}

// marshalValue рекурсивно кодирует значение в буфер.
func (e *Encoder) marshalValue(v reflect.Value) error {
	if !v.IsValid() {
		return fmt.Errorf("invalid or unaddressable value")
	}

	t := v.Type()

	switch t.Kind() {
	case reflect.Bool:
		return e.marshalBool(v.Interface())

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return e.marshalInt(t, v)

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return e.marshalUint(t, v)

	case reflect.Float32, reflect.Float64:
		return e.marshalFloat(t, v)

	case reflect.String:
		return e.marshalString(v)

	case reflect.Ptr:
		return e.marshalPtr(v)

	case reflect.Slice:
		return e.marshalSlice(v)

	case reflect.Array:
		return e.marshalArray(v)

	case reflect.Struct:
		return e.marshalStruct(v)

	default:
		return handleUnsupportedTypeMarshal(t)
	}
}

// marshalBool кодирует булевые типы.
func (e *Encoder) marshalBool(v any) error {
	switch val := v.(type) {
	case bool:
		if val {
			e.buf = append(e.buf, 1)
		} else {
			e.buf = append(e.buf, 0)
		}
	default:
		return fmt.Errorf("unsupported bool type")
	}

	return nil
}

// marshalInt кодирует знаковые целочисленные типы и сохраняет их размер.
func (e *Encoder) marshalInt(t reflect.Type, v reflect.Value) error {
	switch t.Kind() {
	case reflect.Int:
		val := v.Int()
		e.buf = append(e.buf, encodeUint64(uint64(val))...)

	case reflect.Int8:
		val := v.Int()
		e.buf = append(e.buf, byte(uint8(val)))

	case reflect.Int16:
		val := v.Int()
		e.buf = append(e.buf, encodeUint16(uint16(val))...)

	case reflect.Int32:
		val := v.Int()
		e.buf = append(e.buf, encodeUint32(uint32(val))...)

	case reflect.Int64:
		e.buf = append(e.buf, encodeInt64(v.Int())...)

	default:
		return fmt.Errorf("unsupported int type %s", t.String())
	}

	return nil
}

// marshalUint кодирует беззнаковые целочисленные типы и сохраняет их размер.
func (e *Encoder) marshalUint(t reflect.Type, v reflect.Value) error {
	switch t.Kind() {
	case reflect.Uint:
		e.buf = append(e.buf, encodeUint64(uint64(v.Uint()))...)

	case reflect.Uint8:
		e.buf = append(e.buf, byte(v.Uint()))

	case reflect.Uint16:
		e.buf = append(e.buf, encodeUint16(uint16(v.Uint()))...)

	case reflect.Uint32:
		e.buf = append(e.buf, encodeUint32(uint32(v.Uint()))...)

	case reflect.Uint64:
		e.buf = append(e.buf, encodeUint64(v.Uint())...)

	default:
		return fmt.Errorf("unsupported uint type %s", t.String())
	}

	return nil
}

// marshalFloat кодирует плавающие типы.
func (e *Encoder) marshalFloat(t reflect.Type, v reflect.Value) error {
	switch t.Kind() {
	case reflect.Float32:
		e.buf = append(e.buf, encodeFloat32(float32(v.Float()))...)

	case reflect.Float64:
		e.buf = append(e.buf, encodeFloat64(v.Float())...)

	default:
		return fmt.Errorf("unsupported float type %s", t.String())
	}

	return nil
}

// marshalString кодирует строки с префиксом uint32 длины.
func (e *Encoder) marshalString(v reflect.Value) error {
	s := v.String()

	// Validate string length before encoding
	if uint32(len(s)) > maxStringLength {
		return fmt.Errorf("string length %d exceeds maximum", len(s))
	}

	e.buf = append(e.buf, encodeUint32(uint32(len(s)))...)
	e.buf = append(e.buf, []byte(s)...)

	return nil
}

func (e *Encoder) marshalPtr(v reflect.Value) error {
	if v.IsNil() {
		// A nil pointer to a slice marshals the same way Marshal's top-level
		// nil-slice-pointer handling does (as an empty slice), so struct
		// fields of type *[]T behave consistently whether or not they went
		// through Marshal directly.
		if v.Type().Elem().Kind() == reflect.Slice {
			return e.marshalValue(reflect.MakeSlice(v.Type().Elem(), 0, 0))
		}
		return fmt.Errorf("cannot marshal nil value")
	}
	// Delegate to the generic encoder for whatever the pointer points to:
	// a scalar, string, slice, array, struct, or another pointer. Previously
	// this only supported scalar/string element types, so a struct field
	// such as *[]int8 or *SomeStruct would fail with "unsupported pointer
	// element type" even when non-nil.
	return e.marshalValue(v.Elem())
}

// marshalSlice кодирует слайсы с префиксом uint32 количества элементов и элементами.
func (e *Encoder) marshalSlice(v reflect.Value) error {
	if !v.IsValid() || v.Kind() != reflect.Slice {
		return fmt.Errorf("expected slice type, got %s", v.Type().String())
	}

	// Check for nil slice vs empty slice.
	// NOTE: the previous check here was `if v.Interface() == nil`, which is
	// never true for a typed nil slice (a nil []int8 boxed into interface{}
	// is not == nil - classic Go "typed nil" gotcha), so nil slices were
	// silently encoded the same as empty ones. reflect.Value.IsNil() is the
	// correct way to detect this.
	if v.IsNil() {
		e.buf = append(e.buf, encodeUint32(nilSliceSentinel)...)
		return nil
	}

	vLen := v.Len()
	if uint32(vLen) > maxSliceLength {
		return fmt.Errorf("slice too long: %d elements", vLen)
	}

	e.buf = append(e.buf, encodeUint32(uint32(vLen))...)

	for i := 0; i < vLen; i++ {
		if err := e.marshalValue(v.Index(i)); err != nil {
			return fmt.Errorf("marshaling slice element %d: %w", i, err)
		}
	}

	return nil
}

// marshalArray кодирует массивы с элементами (размер в типе).
func (e *Encoder) marshalArray(v reflect.Value) error {
	vLen := int(v.Len())

	// Encode array length prefix for round-trip compatibility
	e.buf = append(e.buf, encodeUint32(uint32(vLen))...)

	for i := 0; i < vLen; i++ {
		elem := v.Index(i)
		if err := e.marshalValue(elem); err != nil {
			return fmt.Errorf("marshaling array element %d: %w", i, err)
		}
	}

	return nil
}

// marshalStruct сериализует поля структуры в порядке их объявления.
func (e *Encoder) marshalStruct(v reflect.Value) error {
	t := v.Type()

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)

		// Skip unexported fields
		if field.PkgPath != "" {
			continue
		}

		// Skip fields with bin:"-" tag
		tag := field.Tag.Get("bin")
		if tag == "-" {
			continue
		}

		fieldVal := v.Field(i)
		if err := e.marshalValue(fieldVal); err != nil {
			return fmt.Errorf("marshaling struct field %s: %w", field.Name, err)
		}
	}

	return nil
}

// handleUnsupportedTypeMarshal возвращает ошибку для неподдерживаемых типов во время сериализации.
func handleUnsupportedTypeMarshal(t reflect.Type) error {
	switch t.Kind() {
	case reflect.Map, reflect.Chan, reflect.Func, reflect.UnsafePointer, reflect.Interface, reflect.Ptr: // ADD reflect.Ptr
		return fmt.Errorf("unsupported type %s (kind: %s) for marshaling",
			t.String(), t.Kind().String())

	default:
		return fmt.Errorf("unsupported type %s (kind: %s) for marshaling",
			t.String(), t.Kind().String())
	}
}

// Helper functions for binary encoding (little-endian).

func encodeInt64(v int64) []byte {
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, uint64(v))
	return b
}

func encodeUint64(v uint64) []byte {
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, v)
	return b
}

func encodeUint32(v uint32) []byte {
	b := make([]byte, 4)
	binary.LittleEndian.PutUint32(b, v)
	return b
}

func encodeUint16(v uint16) []byte {
	b := make([]byte, 2)
	binary.LittleEndian.PutUint16(b, v)
	return b
}

func encodeFloat32(v float32) []byte {
	b := make([]byte, 4)
	binary.LittleEndian.PutUint32(b, math.Float32bits(v))
	return b
}

func encodeFloat64(v float64) []byte {
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, math.Float64bits(v))
	return b
}

// Unmarshal десериализует бинарные данные в Go-значение.
func Unmarshal(data []byte, v any) error {
	if v == nil {
		return fmt.Errorf("cannot unmarshal to nil destination")
	}

	vv := reflect.ValueOf(v)

	// Handle interface containing a concrete value (e.g., any(42))
	if vv.Kind() == reflect.Ptr && vv.Elem().Kind() == reflect.Interface && vv.Elem().Elem().IsValid() {
		return Unmarshal(data, vv.Elem().Elem().Interface())
	}

	if vv.Kind() != reflect.Ptr {
		return fmt.Errorf("destination must be a pointer")
	}

	// FIX: Handle nil pointers - initialize them with new instances before unmarshaling
	if vv.IsNil() {
		tElem := vv.Type().Elem() // Changed from v.Type() to vv.Type()
		vv = reflect.New(tElem)
	}

	decoder := newDecoder(data)
	err := decoder.unmarshalValue(vv.Elem()) // Now vv.Elem() is always valid
	if err != nil {
		return err
	}

	if decoder.offset > len(data) {
		return fmt.Errorf("%w: consumed %d of %d bytes", ErrTruncatedData, decoder.offset, len(data))
	} else if decoder.offset < len(data) {
		return fmt.Errorf("%w: consumed %d of %d bytes", ErrExtraBytes, decoder.offset, len(data))
	}

	return nil
}

// Decoder читает и декодирует значения из бинарных данных.
type Decoder struct {
	data   []byte
	offset int
}

func newDecoder(data []byte) *Decoder {
	return &Decoder{data: data, offset: 0}
}

// unmarshalValue рекурсивно декодирует значение из буфера.
func (d *Decoder) unmarshalValue(v reflect.Value) error {
	if !v.IsValid() || !v.CanAddr() {
		return fmt.Errorf("invalid or non-addressable value")
	}

	// Skip CanSet check for interface kinds - they're handled by extracting concrete types above
	if v.Kind() != reflect.Interface && !v.CanSet() {
		return fmt.Errorf("invalid or non-addressable value")
	}

	t := v.Type()

	switch t.Kind() {
	case reflect.Bool:
		return d.unmarshalBool(v)

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return d.unmarshalInt(t, v)

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return d.unmarshalUint(t, v)

	case reflect.Float32, reflect.Float64:
		return d.unmarshalFloat(t, v)

	case reflect.String:
		return d.unmarshalString(v)

	case reflect.Slice:
		return d.unmarshalSlice(v)

	case reflect.Array:
		return d.unmarshalArray(v)

	case reflect.Struct:
		return d.unmarshalStruct(v)

	case reflect.Ptr:
		return d.unmarshalPtr(v)

	case reflect.Interface:
		return ErrInterfaceTypeUnknown

	default:
		return handleUnsupportedTypeUnmarshal(t)
	}
}

// unmarshalBool декодирует булево значение.
func (d *Decoder) unmarshalBool(v reflect.Value) error {
	if d.offset+1 > len(d.data) {
		return ErrTruncatedData
	}

	val := int8(d.data[d.offset])
	d.offset++

	// Validate boolean value must be 0 or 1
	if val != 0 && val != 1 {
		return ErrInvalidBoolValue
	}

	v.SetBool(val == 1)
	return nil
}

// unmarshalInt декодирует знаковые целочисленные значения.
func (d *Decoder) unmarshalInt(t reflect.Type, v reflect.Value) error {
	switch t.Kind() {
	case reflect.Int:
		if d.offset+8 > len(d.data) {
			return ErrTruncatedData
		}
		val := binary.LittleEndian.Uint64(d.data[d.offset : d.offset+8])
		d.offset += 8
		v.SetInt(int64(val))

	case reflect.Int8:
		if d.offset+1 > len(d.data) {
			return ErrTruncatedData
		}
		val := int8(d.data[d.offset])
		d.offset++
		v.SetInt(int64(val))

	case reflect.Int16:
		if d.offset+2 > len(d.data) {
			return ErrTruncatedData
		}
		val := binary.LittleEndian.Uint16(d.data[d.offset : d.offset+2])
		d.offset += 2
		v.SetInt(int64(int16(val)))

	case reflect.Int32:
		if d.offset+4 > len(d.data) {
			return ErrTruncatedData
		}
		val := binary.LittleEndian.Uint32(d.data[d.offset : d.offset+4])
		d.offset += 4
		v.SetInt(int64(int32(val)))

	case reflect.Int64:
		if d.offset+8 > len(d.data) {
			return ErrTruncatedData
		}
		val := binary.LittleEndian.Uint64(d.data[d.offset : d.offset+8])
		d.offset += 8
		v.SetInt(int64(val))

	default:
		return fmt.Errorf("unsupported int type %s", t.String())
	}

	return nil
}

// unmarshalUint декодирует беззнаковые целочисленные значения.
func (d *Decoder) unmarshalUint(t reflect.Type, v reflect.Value) error {
	switch t.Kind() {
	case reflect.Uint:
		if d.offset+8 > len(d.data) {
			return ErrTruncatedData
		}
		val := binary.LittleEndian.Uint64(d.data[d.offset : d.offset+8])
		d.offset += 8
		v.SetUint(val)

	case reflect.Uint8:
		if d.offset+1 > len(d.data) {
			return ErrTruncatedData
		}
		val := uint8(d.data[d.offset])
		d.offset++
		v.SetUint(uint64(val))

	case reflect.Uint16:
		if d.offset+2 > len(d.data) {
			return ErrTruncatedData
		}
		val := binary.LittleEndian.Uint16(d.data[d.offset : d.offset+2])
		d.offset += 2
		v.SetUint(uint64(val))

	case reflect.Uint32:
		if d.offset+4 > len(d.data) {
			return ErrTruncatedData
		}
		val := binary.LittleEndian.Uint32(d.data[d.offset : d.offset+4])
		d.offset += 4
		v.SetUint(uint64(val))

	case reflect.Uint64:
		if d.offset+8 > len(d.data) {
			return ErrTruncatedData
		}
		val := binary.LittleEndian.Uint64(d.data[d.offset : d.offset+8])
		d.offset += 8
		v.SetUint(val)

	default:
		return fmt.Errorf("unsupported uint type %s", t.String())
	}

	return nil
}

// unmarshalFloat декодирует плавающие значения.
func (d *Decoder) unmarshalFloat(t reflect.Type, v reflect.Value) error {
	switch t.Kind() {
	case reflect.Float32:
		if d.offset+4 > len(d.data) {
			return ErrTruncatedData
		}
		val := binary.LittleEndian.Uint32(d.data[d.offset : d.offset+4])
		d.offset += 4

		f32Bits := uint32(val)
		v.SetFloat(float64(math.Float32frombits(f32Bits)))

	case reflect.Float64:
		if d.offset+8 > len(d.data) {
			return ErrTruncatedData
		}
		val := binary.LittleEndian.Uint64(d.data[d.offset : d.offset+8])
		d.offset += 8

		v.SetFloat(math.Float64frombits(val))

	default:
		return fmt.Errorf("unsupported float type %s", t.String())
	}

	return nil
}

// unmarshalString декодирует UTF-8 строку с префиксом длины.
func (d *Decoder) unmarshalString(v reflect.Value) error {
	if d.offset+4 > len(d.data) {
		return ErrTruncatedData
	}

	strLen := binary.LittleEndian.Uint32(d.data[d.offset : d.offset+4])
	d.offset += 4

	if strLen > maxSliceLength {
		return fmt.Errorf("string length %d exceeds maximum", strLen)
	}

	if d.offset+int(strLen) > len(d.data) {
		return ErrTruncatedData
	}

	str := string(d.data[d.offset : d.offset+int(strLen)])
	d.offset += int(strLen)

	v.SetString(str)
	return nil
}

func (d *Decoder) unmarshalSlice(v reflect.Value) error {
	if d.offset+4 > len(d.data) {
		return ErrTruncatedData
	}

	count := binary.LittleEndian.Uint32(d.data[d.offset : d.offset+4])
	d.offset += 4

	// Handle sentinel value (nil slice): preserve nil state, don't create empty slice.
	if count == nilSliceSentinel {
		v.Set(reflect.Zero(v.Type()))
		return nil
	}

	// Normal slice decoding path - validate length first
	if uint32(count) > maxSliceLength {
		return fmt.Errorf("slice length %d exceeds maximum", count)
	}

	countInt := int(count)

	sliceVal := reflect.MakeSlice(v.Type(), countInt, countInt)

	for i := 0; i < countInt; i++ {
		elemIdx := sliceVal.Index(i)
		if !elemIdx.CanAddr() {
			return fmt.Errorf("slice element %d is not addressable", i)
		}

		if err := d.unmarshalValue(elemIdx); err != nil {
			return fmt.Errorf("decoding slice element %d: %w", i, err)
		}
	}

	v.Set(sliceVal)
	return nil
}

// unmarshalArray декодирует массив (размер в типе, но префикс длины читается).
func (d *Decoder) unmarshalArray(v reflect.Value) error {
	if !v.IsValid() || v.Kind() != reflect.Array {
		return fmt.Errorf("expected array type, got %s", v.Type().String())
	}

	vLen := int(v.Len())

	// Read length prefix for round-trip compatibility (matches marshal behavior)
	if d.offset+4 > len(d.data) {
		return ErrTruncatedData
	}

	lengthPrefix := binary.LittleEndian.Uint32(d.data[d.offset : d.offset+4])
	d.offset += 4

	// Verify the length prefix matches expected array size
	if lengthPrefix != uint32(vLen) {
		return fmt.Errorf("array length mismatch: got %d, want %d", lengthPrefix, vLen)
	}

	for i := 0; i < vLen; i++ {
		elemVal := v.Index(i)
		if !elemVal.CanAddr() {
			return fmt.Errorf("array element %d is not addressable", i)
		}

		if err := d.unmarshalValue(elemVal); err != nil {
			return fmt.Errorf("decoding array element %d: %w", i, err)
		}
	}

	return nil
}

func (d *Decoder) unmarshalStruct(v reflect.Value) error {
	t := v.Type()

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)

		// Skip unexported fields
		if field.PkgPath != "" {
			continue
		}

		// Skip fields with bin:"-" tag
		tag := field.Tag.Get("bin")
		if tag == "-" {
			continue
		}

		fieldVal := v.Field(i)

		if err := d.unmarshalValue(fieldVal); err != nil {
			return fmt.Errorf("decoding struct field %s: %w", field.Name, err)
		}

	}

	return nil
}

func (d *Decoder) unmarshalPtr(v reflect.Value) error {
	// Allocate the pointee if the pointer is nil, so we always have a
	// concrete, addressable value to decode into.
	if v.IsNil() {
		v.Set(reflect.New(v.Type().Elem()))
	}

	elem := v.Elem()

	if elem.Kind() == reflect.Interface {
		return ErrInterfaceTypeUnknown
	}

	// Delegate to the generic decoder for whatever the pointer points to:
	// a scalar, string, slice, array, struct, or another pointer (handling
	// **T by recursing here again). unmarshalValue's own dispatch already
	// knows how to handle every one of these kinds correctly, so there is
	// no need to special-case them here - doing so previously meant only
	// scalar pointee types worked, and *[]int8/*[]uint8 (or any *struct,
	// *array, *string) fell into a default branch and errored with
	// "cannot unmarshal into nil ...".
	return d.unmarshalValue(elem)
}

// handleUnsupportedTypeUnmarshal возвращает ошибку для неподдерживаемых типов во время десериализации.
func handleUnsupportedTypeUnmarshal(t reflect.Type) error {
	return fmt.Errorf("unsupported type %s (kind: %s) for unmarshaling",
		t.String(), t.Kind().String())
}
