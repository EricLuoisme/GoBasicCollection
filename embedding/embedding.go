package main

import "fmt"

// Position 例子是Struct Embedding, 通过struct之间直接引入另一个struct作为anonymous的属性, 可以直接调用它支持的方法
type Position struct {
	x float64
	y float64
}

// SpecialPosition 这里一定要用Position而不是Position的Pointer, 否则golang无法初始化 (记住struct最底层的一个依赖, 一定要是实体才可以)
type SpecialPosition struct {
	Position
}

func (sp *SpecialPosition) MoveSpecial(x, y float64) {
	sp.x += x * x
	sp.y += y * y
}

func (p *Position) Move(x, y float64) {
	p.x += x
	p.y += y
}

func (p *Position) Teleport(x, y float64) {
	p.x = x
	p.y = y
}

type Player struct {
	//posX float64
	//posY float64
	*Position // 直接引用指针, embedding
}

func NewPlayer() *Player {
	return &Player{Position: &Position{}} // golang会帮助初始化float64初始值
}

//func (p *Player) Move(x, y float64) {
//	p.posX += x
//	p.posY += y
//}
//
//func (p *Player) Teleport(x, y float64) {
//	p.posX = x
//	p.posY = y
//}

type Enemy struct {
	*SpecialPosition
}

func NewEnemy() *Enemy {
	return &Enemy{SpecialPosition: &SpecialPosition{}} // golang会帮助初始化float64初始值
}

//func (p *Enemy) Move(x, y float64) {
//	p.posX += x
//	p.posY += y
//}
//
//func (p *Enemy) Teleport(x, y float64) {
//	p.posX = x
//	p.posY = y
//}

func main() {
	// 通过embedding, player这个struct可以完全使用所有实现Position的func
	player := NewPlayer()
	player.Move(0, 1)
	fmt.Printf("%+v\n", player.Position)

	// boss
	boss := NewEnemy()
	boss.Move(1, 2)
	fmt.Printf("%+v\n", boss.Position)
	fmt.Printf("%+v\n", boss.SpecialPosition)

	boss.MoveSpecial(2, 4)
	fmt.Printf("%+v\n", boss.SpecialPosition)
}
