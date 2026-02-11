package collider

type bitmask uint32

// Predefined layer constants
const (
	LayerPlayer     uint32 = 0
	LayerEnemy      uint32 = 1
	LayerProjectile uint32 = 2
	LayerWall       uint32 = 3
	LayerTrigger    uint32 = 4
)

/**
 * Bitmask is a utility type for managing layers in a game.
 * It provides methods for setting, clearing, and checking bits,
 * as well as checking for matches between two bitmasks.
 */

func NewBitmask() bitmask {
	return 0
}

func (b *bitmask) SetBit(bit uint32) {
	*b |= bitmask(1 << bit)
}

func (b *bitmask) ClearBit(bit uint32) {
	*b &^= bitmask(1 << bit)
}

func (b *bitmask) IsSet(bit uint32) bool {
	return (*b & bitmask(1<<bit)) != 0
}

func (b *bitmask) CanMatch(other bitmask) bool {
	return (*b & other) != 0
}

func (b *bitmask) SetLayers(layers ...uint32) {
	for _, layer := range layers {
		b.SetBit(layer)
	}
}

func (b *bitmask) HasAny(other bitmask) bool {
	return (*b & other) != 0
}

func (b *bitmask) HasAll(other bitmask) bool {
	return (*b & other) == other
}
