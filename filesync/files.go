package filesync

type Files []File

func (f Files) SpaceNeeded() int64 {
	var space int64

	for _, v := range f {
		space += v.Source.Info.Size()

		if v.Destination.Info != nil {
			space -= v.Destination.Info.Size()
		}
	}

	return space
}
