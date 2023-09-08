package main

import "fmt"

// 在这个例子中，Text和Image都是Element的实现。我们有两个访问者，RenderVisitor和WordCountVisitor。RenderVisitor负责渲染文档的内容，而WordCountVisitor则负责统计字数。
// 这种方式使得在不修改Text或Image类的情况下，我们可以轻松地为它们添加新的操作，只需添加一个新的访问者类即可。


// Element 接口
type Element interface {
	Accept(visitor Visitor)
}

// Visitor 接口
type Visitor interface {
	VisitText(text *Text)
	VisitImage(image *Image)
}

// Text 类型
type Text struct {
	content string
}

func (t *Text) Accept(visitor Visitor) {
	visitor.VisitText(t)
}

// Image 类型
type Image struct {
	altText string
}

func (i *Image) Accept(visitor Visitor) {
	visitor.VisitImage(i)
}

// RenderVisitor 渲染访问者
type RenderVisitor struct{}

func (r *RenderVisitor) VisitText(text *Text) {
	fmt.Println("Rendering Text:", text.content)
}

func (r *RenderVisitor) VisitImage(image *Image) {
	fmt.Println("Rendering Image with alt:", image.altText)
}

// WordCountVisitor 字数统计访问者
type WordCountVisitor struct{}

func (w *WordCountVisitor) VisitText(text *Text) {
	fmt.Println("Word count for text:", len(text.content))
}

func (w *WordCountVisitor) VisitImage(image *Image) {
	fmt.Println("Word count for image alt text:", len(image.altText))
}

func main() {
	elements := []Element{
		&Text{content: "Hello World"},
		&Image{altText: "An Image"},
	}

	renderVisitor := &RenderVisitor{}
	wordCountVisitor := &WordCountVisitor{}

	for _, el := range elements {
		el.Accept(renderVisitor)
		el.Accept(wordCountVisitor)
	}
}
