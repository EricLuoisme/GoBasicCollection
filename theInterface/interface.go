package main

import (
	"fmt"
	"math/rand"
)

type Player interface {
	KickBall()
}

type FootballPlayer struct {
	stamina int
	power   int
}

// KickBall 属于FootballPlayer的struct可以使用的方法, 由于与interface中的方法相同,
// 实际上就是让FootballPlayer struct实现了Player interface
func (f *FootballPlayer) KickBall() {
	shot := f.stamina + f.power
	fmt.Println("I'm kicking the ball", shot)
}

func main() {

	// 初始化一个array, 长度是11人
	//team := make([]FootballPlayer, 11)

	team := make([]Player, 11, 11) // 相对于使用struct的数组, 使用interface的会更加抽象, 方便解耦
	// 遍历Random初始化
	for i := 0; i < len(team); i++ {
		team[i] = &FootballPlayer{
			stamina: rand.Intn(10),
			power:   rand.Intn(10),
		}
	}

	// 循环调用
	for i := 0; i < len(team); i++ {
		team[i].KickBall()
	}
}
