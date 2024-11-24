package erc20

import (
	"gobius/utils"
	"math"
	"math/big"
	"strconv"
	"strings"

	"hash"
	"hash/fnv"

	"github.com/ethereum/go-ethereum/common"
)

func FNV64(data []byte) int64 {
	algorithm := fnv.New64()
	return uint64Hasher(algorithm, data)
}

func uint64Hasher(algorithm hash.Hash64, data []byte) int64 {
	algorithm.Write(data)
	return int64(algorithm.Sum64())
}

type TokenERC20 struct {
	//chainId  i
	Address  common.Address
	Decimals int64
	Symbol   string
	Name     string
	Id       int64 `json:"-"`
}

func NewTokenERC20(address common.Address, decimals int64, symbol string, name string) *TokenERC20 {

	Id := FNV64(address.Bytes())

	return &TokenERC20{
		//chainId:  chainI
		Address:  address,
		Decimals: decimals,
		Symbol:   symbol,
		Name:     name,
		Id:       Id,
	}
}

var _ten = big.NewInt(10)

func (t *TokenERC20) DecimalExp() *big.Int {
	return new(big.Int).Exp(_ten, big.NewInt(t.Decimals), nil)
}

func (t *TokenERC20) Init() {
	t.Id = FNV64(t.Address.Bytes())
}

func (t *TokenERC20) ID() int64 {
	if t.Id == 0 {
		t.Id = FNV64(t.Address.Bytes())
	}
	return t.Id
}

func (t *TokenERC20) SortsBefore(other *TokenERC20) bool {
	//return strings.ToLower(t.Address.String()) < strings.ToLower(other.Address.String())
	token0Rep := big.NewInt(0).SetBytes(t.Address.Bytes())
	token1Rep := big.NewInt(0).SetBytes(other.Address.Bytes())

	return token0Rep.Cmp(token1Rep) < 0

}

func (t *TokenERC20) Equals(other *TokenERC20) bool {
	return t.Address == other.Address
}

func (t *TokenERC20) EqualsFast(other *TokenERC20) bool {
	return t.Id == other.Id
}

func (t *TokenERC20) ONE() *big.Int {
	return new(big.Int).Exp(_ten, big.NewInt(t.Decimals), nil)
}

func (t *TokenERC20) OfAmount(amount int64, scale int64) *big.Int {
	value := new(big.Int).Mul(t.ONE(), big.NewInt(amount))
	if scale != 1 {
		value.Div(value, big.NewInt(scale))
	}
	return value
}

// convert big.Int amount of Token into float64 value accounting for it's decimals
func (t *TokenERC20) ToFloat(amount *big.Int) float64 {
	if amount == nil {
		return math.NaN()
	}
	var base, exponent = big.NewInt(10), big.NewInt(int64(t.Decimals))
	denominator := base.Exp(base, exponent, nil)
	tokensSentFloat := new(big.Float).SetInt(amount)
	denominatorFloat := new(big.Float).SetInt(denominator)
	final, _ := new(big.Float).Quo(tokensSentFloat, denominatorFloat).Float64()
	return final
}

func (t *TokenERC20) FormatFixed(amount *big.Int) string {
	val, _ := utils.FormatFixed(amount, int(t.Decimals))

	return val
}

func (t *TokenERC20) FromFloat(amount float64) *big.Int {
	var base, exponent = big.NewInt(10), big.NewInt(int64(t.Decimals))
	multiplier := base.Exp(base, exponent, nil)

	bigval := new(big.Float)
	bigval.SetFloat64(amount)
	// Set precision if required.
	bigval.SetPrec(256)

	coin := new(big.Float)
	coin.SetInt(multiplier)

	bigval.Mul(bigval, coin)

	result := new(big.Int)
	bigval.Int(result) // store converted number in result

	return result
}

func (t *TokenERC20) StringToBigInt(value string) *big.Int {

	valueAsFloat, err := strconv.ParseFloat(strings.TrimSpace(value), 64)

	if err != nil {
		return nil
	}

	return t.FromFloat(valueAsFloat)

}
