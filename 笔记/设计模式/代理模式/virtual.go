type VirtualSubject interface {
    Display()
}

type RealImage struct {
    filename string
}

func NewRealImage(filename string) *RealImage {
    // Here we can simulate loading the image from disk.
    fmt.Printf("Loading image %s from disk.\n", filename)
    return &RealImage{filename: filename}
}

func (ri *RealImage) Display() {
    fmt.Printf("Displaying image %s.\n", ri.filename)
}

type VirtualProxy struct {
    filename  string
    realImage *RealImage
}

func (vp *VirtualProxy) Display() {
    if vp.realImage == nil {
        vp.realImage = NewRealImage(vp.filename)
    }
    vp.realImage.Display()
}
