package dependencyinjection

import "fmt"

/**
依赖注入核心: 在golang里面, 对于struct尽可能的依赖interface, 而不是实际的struct, 来达到更动态的依赖注入
*/

// SaftyPlacer 从struct转换为interface
type SaftyPlacer interface {
	placeSafeties()
}

// IceSaftyPlacer 可以看作是需要实现SaftyPlacer的一种, 由于需要进行具体操作, 这里先设置为struct
type IceSaftyPlacer struct {
	//
}

// IceSaftyPlacer 需要实现SaftyPlacer接口, 所以进行方法实现
func (isp IceSaftyPlacer) PlaceSafeties() {
	fmt.Println("placing my safeties...")
}

type RockClimber struct {
	kind         int
	rocksClimbed int
	sp           SaftyPlacer // 相对于使用struct, 将SaftyPlacer抽取为interface可以实现依赖注入
}

// RockClimber 的生成器(类似构造方法)
func newRockClimber(sp SaftyPlacer) *RockClimber {
	return &RockClimber{
		sp: sp,
	}
}

func (rc *RockClimber) climbRock() {
	rc.rocksClimbed++
	if rc.rocksClimbed == 10 {
		rc.sp.placeSafeties()
	}
}
