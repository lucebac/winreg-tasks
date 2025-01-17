// Code generated by kaitai-struct-compiler from a .ksy source file. DO NOT EDIT.

package generated

import "github.com/kaitai-io/kaitai_struct_go_runtime/kaitai"

type DynamicInfo struct {
	Magic []byte
	CreationTime uint64
	LastRunTime uint64
	TaskState uint32
	LastErrorCode uint32
	LastSuccessfulRunTime uint64
	_io *kaitai.Stream
	_root *DynamicInfo
	_parent interface{}
}
func NewDynamicInfo() *DynamicInfo {
	return &DynamicInfo{
	}
}

func (this *DynamicInfo) Read(io *kaitai.Stream, parent interface{}, root *DynamicInfo) (err error) {
	this._io = io
	this._parent = parent
	this._root = root

	tmp1, err := this._io.ReadBytes(int(4))
	if err != nil {
		return err
	}
	tmp1 = tmp1
	this.Magic = tmp1
	tmp2, err := this._io.ReadU8le()
	if err != nil {
		return err
	}
	this.CreationTime = uint64(tmp2)
	tmp3, err := this._io.ReadU8le()
	if err != nil {
		return err
	}
	this.LastRunTime = uint64(tmp3)
	tmp4, err := this._io.ReadU4le()
	if err != nil {
		return err
	}
	this.TaskState = uint32(tmp4)
	tmp5, err := this._io.ReadU4le()
	if err != nil {
		return err
	}
	this.LastErrorCode = uint32(tmp5)
	tmp6, err := this._io.EOF()
	if err != nil {
		return err
	}
	if (!(tmp6)) {
		tmp7, err := this._io.ReadU8le()
		if err != nil {
			return err
		}
		this.LastSuccessfulRunTime = uint64(tmp7)
	}
	return err
}
