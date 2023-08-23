package filesync

type Files []File

func (f Files) SpaceNeeded() int64 {
	var space int64

	for _, v := range f {
		space += v.Source.Size()

		if v.Destination != nil {
			space -= v.Destination.Size()
		}
	}

	return space
}
