package wallpaper_test

import (
	. "gopkg.in/check.v1"
	"math"
	"math/cmplx"
	"wallpaper/entities/formula"
	"wallpaper/entities/formula/wallpaper"
	"wallpaper/entities/utility"
)

type HexagonalWallpaper struct {
	newFormula *wallpaper.Formula
}

var _ = Suite(&HexagonalWallpaper{})

func (suite *HexagonalWallpaper) SetUpTest (checker *C) {
	suite.newFormula = &wallpaper.Formula{
		LatticeType:     wallpaper.Hexagonal,
		LatticeSize:     &wallpaper.Dimensions{
			Width:  2,
			Height: -0.5,
		},
		Lattice:         nil,
		Multiplier:      complex(1, 0),
		WavePackets:     []*wallpaper.WavePacket{
			{
				Terms: []*formula.EisensteinFormulaTerm{
					{
						PowerN: 1,
						PowerM: -2,
					},
				},
				Multiplier: complex(1, 0),
			},
		},
	}

	suite.newFormula.Setup()
}

func (suite *HexagonalWallpaper) TestSetupCreatesLatticeVectors (checker *C) {
	checker.Assert(real(suite.newFormula.Lattice.XLatticeVector), utility.NumericallyCloseEnough{}, 1, 1e-6)
	checker.Assert(imag(suite.newFormula.Lattice.XLatticeVector), utility.NumericallyCloseEnough{}, 0, 1e-6)

	checker.Assert(real(suite.newFormula.Lattice.YLatticeVector), utility.NumericallyCloseEnough{}, -0.5, 1e-6)
	checker.Assert(imag(suite.newFormula.Lattice.YLatticeVector), utility.NumericallyCloseEnough{}, math.Sqrt(3.0)/2.0, 1e-6)
}

func (suite *HexagonalWallpaper) TestSetupAddsLockedPairs (checker *C) {
	checker.Assert(suite.newFormula.WavePackets[0].Terms, HasLen, 3)
	checker.Assert(suite.newFormula.WavePackets[0].Terms[1].PowerN, Equals, -2)
	checker.Assert(suite.newFormula.WavePackets[0].Terms[1].PowerM, Equals, 1)
	checker.Assert(suite.newFormula.WavePackets[0].Terms[2].PowerN, Equals, 1)
	checker.Assert(suite.newFormula.WavePackets[0].Terms[2].PowerM, Equals, 1)
}

func (suite *HexagonalWallpaper) TestCalculationOfPoints (checker *C) {
	calculation := suite.newFormula.Calculate(complex(math.Sqrt(3), -1 * math.Sqrt(3)))
	total := calculation.Total

	expectedAnswer := (cmplx.Exp(complex(0, 2 * math.Pi * (3 + math.Sqrt(3)))) +
		cmplx.Exp(complex(0, 2 * math.Pi * (-2 * math.Sqrt(3)))) +
		cmplx.Exp(complex(0, 2 * math.Pi * (-3 + math.Sqrt(3))))) / 3

	checker.Assert(real(total), utility.NumericallyCloseEnough{}, real(expectedAnswer), 1e-6)
	checker.Assert(imag(total), utility.NumericallyCloseEnough{}, imag(expectedAnswer), 1e-6)
}

type HexagonalWallpaperHasSymmetryTest struct {
	baseWavePacket *wallpaper.WavePacket
	wallpaperMultiplier complex128
}

var _ = Suite(&HexagonalWallpaperHasSymmetryTest{})

func (suite *HexagonalWallpaperHasSymmetryTest) SetUpTest(checker *C) {
	suite.baseWavePacket = &wallpaper.WavePacket{
		Terms:[]*formula.EisensteinFormulaTerm{
			{
				PowerN:         8,
				PowerM:         -3,
			},
		},
		Multiplier: complex(1, 0),
	}

	suite.wallpaperMultiplier = complex(1, 0)
}

func (suite *HexagonalWallpaperHasSymmetryTest) TestHexagonalWillAlwaysHaveP3(checker *C) {
	newFormula := wallpaper.Formula{
		LatticeType:     wallpaper.Hexagonal,
		LatticeSize:     nil,
		Lattice:         nil,
		Multiplier:      complex(2, 0),
		WavePackets:     []*wallpaper.WavePacket{
			suite.baseWavePacket,
		},
	}
	err := newFormula.Setup()
	checker.Assert(err, IsNil)

	checker.Assert(newFormula.WavePackets, HasLen, 1)
	checker.Assert(newFormula.HasSymmetry(wallpaper.P3), Equals, true)
	checker.Assert(newFormula.HasSymmetry(wallpaper.P31m), Equals, false)
	checker.Assert(newFormula.HasSymmetry(wallpaper.P3m1), Equals, false)
	checker.Assert(newFormula.HasSymmetry(wallpaper.P6), Equals, false)
	checker.Assert(newFormula.HasSymmetry(wallpaper.P6m), Equals, false)
}

func (suite *HexagonalWallpaperHasSymmetryTest) TestHexagonalMayHaveSymmetryForP31m(checker *C) {
	newFormula := wallpaper.Formula{
		LatticeType:     wallpaper.Hexagonal,
		LatticeSize:     nil,
		Lattice:         nil,
		Multiplier:      complex(2, 0),
		WavePackets:     []*wallpaper.WavePacket{
			suite.baseWavePacket,
			{
				Terms: []*formula.EisensteinFormulaTerm{
					{
						PowerN:         suite.baseWavePacket.Terms[0].PowerM * -1,
						PowerM:         suite.baseWavePacket.Terms[0].PowerN * -1,
					},
				},
				Multiplier: suite.baseWavePacket.Multiplier,
			},
		},
	}
	err := newFormula.Setup()
	checker.Assert(err, IsNil)

	checker.Assert(newFormula.HasSymmetry(wallpaper.P3), Equals, true)
	checker.Assert(newFormula.HasSymmetry(wallpaper.P31m), Equals, false)
	checker.Assert(newFormula.HasSymmetry(wallpaper.P3m1), Equals, true)
	checker.Assert(newFormula.HasSymmetry(wallpaper.P6), Equals, false)
	checker.Assert(newFormula.HasSymmetry(wallpaper.P6m), Equals, false)
}

func (suite *HexagonalWallpaperHasSymmetryTest) TestHexagonalMayHaveSymmetryForP3m1(checker *C) {
	newFormula := wallpaper.Formula{
		LatticeType:     wallpaper.Hexagonal,
		LatticeSize:     nil,
		Lattice:         nil,
		Multiplier:      complex(2, 0),
		WavePackets:     []*wallpaper.WavePacket{
			suite.baseWavePacket,
			{
				Terms: []*formula.EisensteinFormulaTerm{
					{
						PowerN:         suite.baseWavePacket.Terms[0].PowerM * -1,
						PowerM:         suite.baseWavePacket.Terms[0].PowerN * -1,
					},
				},
				Multiplier: complex(1, 0),
			},
		},
	}
	err := newFormula.Setup()
	checker.Assert(err, IsNil)

	checker.Assert(newFormula.HasSymmetry(wallpaper.P3), Equals, true)
	checker.Assert(newFormula.HasSymmetry(wallpaper.P31m), Equals, false)
	checker.Assert(newFormula.HasSymmetry(wallpaper.P3m1), Equals, true)
	checker.Assert(newFormula.HasSymmetry(wallpaper.P6), Equals, false)
	checker.Assert(newFormula.HasSymmetry(wallpaper.P6m), Equals, false)
}

func (suite *HexagonalWallpaperHasSymmetryTest) TestHexagonalMayHaveSymmetryForP6(checker *C) {
	newFormula := wallpaper.Formula{
		LatticeType:     wallpaper.Hexagonal,
		LatticeSize:     nil,
		Lattice:         nil,
		Multiplier:      complex(2, 0),
		WavePackets:     []*wallpaper.WavePacket{
			suite.baseWavePacket,
			{
				Terms: []*formula.EisensteinFormulaTerm{
					{
						PowerN:         suite.baseWavePacket.Terms[0].PowerN * -1,
						PowerM:         suite.baseWavePacket.Terms[0].PowerM * -1,
					},
				},
				Multiplier: complex(1, 0),
			},
		},
	}
	err := newFormula.Setup()
	checker.Assert(err, IsNil)

	checker.Assert(newFormula.HasSymmetry(wallpaper.P3), Equals, true)
	checker.Assert(newFormula.HasSymmetry(wallpaper.P31m), Equals, false)
	checker.Assert(newFormula.HasSymmetry(wallpaper.P3m1), Equals, false)
	checker.Assert(newFormula.HasSymmetry(wallpaper.P6), Equals, true)
	checker.Assert(newFormula.HasSymmetry(wallpaper.P6m), Equals, false)
}

func (suite *HexagonalWallpaperHasSymmetryTest) TestHexagonalMayHaveSymmetryForP6m(checker *C) {
	newFormula := wallpaper.Formula{
		LatticeType:     wallpaper.Hexagonal,
		LatticeSize:     nil,
		Lattice:         nil,
		Multiplier:      complex(2, 0),
		WavePackets:     []*wallpaper.WavePacket{
			suite.baseWavePacket,
			{
				Terms: []*formula.EisensteinFormulaTerm{
					{
						PowerN:         suite.baseWavePacket.Terms[0].PowerN * -1,
						PowerM:         suite.baseWavePacket.Terms[0].PowerM * -1,
					},
				},
				Multiplier: suite.baseWavePacket.Multiplier,
			},
			{
				Terms: []*formula.EisensteinFormulaTerm{
					{
						PowerN:         suite.baseWavePacket.Terms[0].PowerM,
						PowerM:         suite.baseWavePacket.Terms[0].PowerN,
					},
				},
				Multiplier: suite.baseWavePacket.Multiplier,
			},
			{
				Terms: []*formula.EisensteinFormulaTerm{
					{
						PowerN:         suite.baseWavePacket.Terms[0].PowerM * -1,
						PowerM:         suite.baseWavePacket.Terms[0].PowerN * -1,
					},
				},
				Multiplier: suite.baseWavePacket.Multiplier,
			},
		},
	}
	err := newFormula.Setup()
	checker.Assert(err, IsNil)

	checker.Assert(newFormula.HasSymmetry(wallpaper.P3), Equals, true)
	checker.Assert(newFormula.HasSymmetry(wallpaper.P31m), Equals, true)
	checker.Assert(newFormula.HasSymmetry(wallpaper.P3m1), Equals, true)
	checker.Assert(newFormula.HasSymmetry(wallpaper.P6), Equals, true)
	checker.Assert(newFormula.HasSymmetry(wallpaper.P6m), Equals, true)
}

type HexagonalCreatedWithDesiredSymmetry struct {
	baseWavePacket *wallpaper.WavePacket
	wallpaperMultiplier complex128
}

var _ = Suite(&HexagonalCreatedWithDesiredSymmetry{})

func (suite *HexagonalCreatedWithDesiredSymmetry) SetUpTest (checker *C) {
	suite.baseWavePacket = &wallpaper.WavePacket{
		Terms:[]*formula.EisensteinFormulaTerm{
			{
				PowerN: 1,
				PowerM: -2,
			},
		},
		Multiplier: complex(1, 0),
	}
	suite.wallpaperMultiplier = complex(1, 0)
}

func (suite *HexagonalCreatedWithDesiredSymmetry) TestCreateWallpaperWithP31m(checker *C) {
	newFormula := wallpaper.Formula{
		LatticeType:     wallpaper.Hexagonal,
		Multiplier:      complex(2, 0),
		WavePackets:     []*wallpaper.WavePacket{
			suite.baseWavePacket,
		},
		DesiredSymmetry: wallpaper.P31m,
	}
	err := newFormula.Setup()
	checker.Assert(err, IsNil)

	checker.Assert(newFormula.WavePackets, HasLen, 2)
	checker.Assert(newFormula.WavePackets[0].Terms, HasLen, 3)

	checker.Assert(newFormula.WavePackets[1].Terms[0].PowerM, Equals, 1)
	checker.Assert(newFormula.WavePackets[1].Terms[0].PowerN, Equals, -2)

	checker.Assert(newFormula.HasSymmetry(wallpaper.P3), Equals, true)
	checker.Assert(newFormula.HasSymmetry(wallpaper.P31m), Equals, true)
}

func (suite *HexagonalCreatedWithDesiredSymmetry) TestCreateWallpaperWithP3m1(checker *C) {
	newFormula := wallpaper.Formula{
		LatticeType:     wallpaper.Hexagonal,
		Multiplier:      complex(2, 0),
		WavePackets:     []*wallpaper.WavePacket{
			suite.baseWavePacket,
		},
		DesiredSymmetry: wallpaper.P3m1,
	}
	err := newFormula.Setup()
	checker.Assert(err, IsNil)
	checker.Assert(newFormula.WavePackets, HasLen, 2)
	checker.Assert(newFormula.WavePackets[0].Terms, HasLen, 3)

	checker.Assert(newFormula.WavePackets[1].Terms[0].PowerM, Equals, -1)
	checker.Assert(newFormula.WavePackets[1].Terms[0].PowerN, Equals, 2)

	checker.Assert(newFormula.HasSymmetry(wallpaper.P3), Equals, true)
	checker.Assert(newFormula.HasSymmetry(wallpaper.P3m1), Equals, true)
}

func (suite *HexagonalCreatedWithDesiredSymmetry) TestCreateWallpaperWithP6(checker *C) {
	newFormula := wallpaper.Formula{
		LatticeType:     wallpaper.Hexagonal,
		Multiplier:      complex(2, 0),
		WavePackets:     []*wallpaper.WavePacket{
			suite.baseWavePacket,
		},
		DesiredSymmetry: wallpaper.P6,
	}
	err := newFormula.Setup()
	checker.Assert(err, IsNil)
	checker.Assert(newFormula.WavePackets, HasLen, 2)
	checker.Assert(newFormula.WavePackets[0].Terms, HasLen, 3)

	checker.Assert(newFormula.WavePackets[1].Terms[0].PowerM, Equals, 2)
	checker.Assert(newFormula.WavePackets[1].Terms[0].PowerN, Equals, -1)

	checker.Assert(newFormula.HasSymmetry(wallpaper.P3), Equals, true)
	checker.Assert(newFormula.HasSymmetry(wallpaper.P6), Equals, true)
}

func (suite *HexagonalCreatedWithDesiredSymmetry) TestCreateWallpaperWithP6m(checker *C) {
	newFormula := wallpaper.Formula{
		LatticeType:     wallpaper.Hexagonal,
		Multiplier:      complex(2, 0),
		WavePackets:     []*wallpaper.WavePacket{
			suite.baseWavePacket,
		},
		DesiredSymmetry: wallpaper.P6m,
	}
	err := newFormula.Setup()
	checker.Assert(err, IsNil)
	checker.Assert(newFormula.WavePackets, HasLen, 4)
	checker.Assert(newFormula.WavePackets[0].Terms, HasLen, 3)

	checker.Assert(newFormula.WavePackets[1].Terms[0].PowerM, Equals, 2)
	checker.Assert(newFormula.WavePackets[1].Terms[0].PowerN, Equals, -1)

	checker.Assert(newFormula.WavePackets[2].Terms[0].PowerM, Equals, 1)
	checker.Assert(newFormula.WavePackets[2].Terms[0].PowerN, Equals, -2)

	checker.Assert(newFormula.WavePackets[3].Terms[0].PowerM, Equals, -1)
	checker.Assert(newFormula.WavePackets[3].Terms[0].PowerN, Equals, 2)

	checker.Assert(newFormula.HasSymmetry(wallpaper.P3), Equals, true)
	checker.Assert(newFormula.HasSymmetry(wallpaper.P6m), Equals, true)
}
