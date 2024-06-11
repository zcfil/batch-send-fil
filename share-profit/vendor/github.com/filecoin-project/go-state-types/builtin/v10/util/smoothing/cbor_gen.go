// Code generated by github.com/whyrusleeping/cbor-gen. DO NOT EDIT.

package smoothing

import (
	"fmt"
	"io"
	"sort"

	cid "github.com/ipfs/go-cid"
	cbg "github.com/whyrusleeping/cbor-gen"
	xerrors "golang.org/x/xerrors"
)

var _ = xerrors.Errorf
var _ = cid.Undef
var _ = sort.Sort

var lengthBufFilterEstimate = []byte{130}

func (t *FilterEstimate) MarshalCBOR(w io.Writer) error {
	if t == nil {
		_, err := w.Write(cbg.CborNull)
		return err
	}
	if _, err := w.Write(lengthBufFilterEstimate); err != nil {
		return err
	}

	// t.PositionEstimate (big.Int) (struct)
	if err := t.PositionEstimate.MarshalCBOR(w); err != nil {
		return err
	}

	// t.VelocityEstimate (big.Int) (struct)
	if err := t.VelocityEstimate.MarshalCBOR(w); err != nil {
		return err
	}
	return nil
}

func (t *FilterEstimate) UnmarshalCBOR(r io.Reader) error {
	*t = FilterEstimate{}

	br := cbg.GetPeeker(r)
	scratch := make([]byte, 8)

	maj, extra, err := cbg.CborReadHeaderBuf(br, scratch)
	if err != nil {
		return err
	}
	if maj != cbg.MajArray {
		return fmt.Errorf("cbor input should be of type array")
	}

	if extra != 2 {
		return fmt.Errorf("cbor input had wrong number of fields")
	}

	// t.PositionEstimate (big.Int) (struct)

	{

		if err := t.PositionEstimate.UnmarshalCBOR(br); err != nil {
			return xerrors.Errorf("unmarshaling t.PositionEstimate: %w", err)
		}

	}
	// t.VelocityEstimate (big.Int) (struct)

	{

		if err := t.VelocityEstimate.UnmarshalCBOR(br); err != nil {
			return xerrors.Errorf("unmarshaling t.VelocityEstimate: %w", err)
		}

	}
	return nil
}
