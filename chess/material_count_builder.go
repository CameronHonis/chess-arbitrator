package chess

type MaterialCountBuilder struct {
	materialCount *MaterialCount
}

func NewMaterialCountBuilder() *MaterialCountBuilder {
	return &MaterialCountBuilder{
		materialCount: &MaterialCount{},
	}
}

func (builder *MaterialCountBuilder) WithWhitePawnCount(count uint8) *MaterialCountBuilder {
	builder.materialCount.WhitePawnCount = count
	return builder
}

func (builder *MaterialCountBuilder) WithWhiteKnightCount(count uint8) *MaterialCountBuilder {
	builder.materialCount.WhiteKnightCount = count
	return builder
}

func (builder *MaterialCountBuilder) WithWhiteLightBishopCount(count uint8) *MaterialCountBuilder {
	builder.materialCount.WhiteLightBishopCount = count
	return builder
}

func (builder *MaterialCountBuilder) WithWhiteDarkBishopCount(count uint8) *MaterialCountBuilder {
	builder.materialCount.WhiteDarkBishopCount = count
	return builder
}

func (builder *MaterialCountBuilder) WithWhiteRookCount(count uint8) *MaterialCountBuilder {
	builder.materialCount.WhiteRookCount = count
	return builder
}

func (builder *MaterialCountBuilder) WithWhiteQueenCount(count uint8) *MaterialCountBuilder {
	builder.materialCount.WhiteQueenCount = count
	return builder
}

func (builder *MaterialCountBuilder) WithBlackPawnCount(count uint8) *MaterialCountBuilder {
	builder.materialCount.BlackPawnCount = count
	return builder
}

func (builder *MaterialCountBuilder) WithBlackKnightCount(count uint8) *MaterialCountBuilder {
	builder.materialCount.BlackKnightCount = count
	return builder
}

func (builder *MaterialCountBuilder) WithBlackLightBishopCount(count uint8) *MaterialCountBuilder {
	builder.materialCount.BlackLightBishopCount = count
	return builder
}

func (builder *MaterialCountBuilder) WithBlackDarkBishopCount(count uint8) *MaterialCountBuilder {
	builder.materialCount.BlackDarkBishopCount = count
	return builder
}

func (builder *MaterialCountBuilder) WithBlackRookCount(count uint8) *MaterialCountBuilder {
	builder.materialCount.BlackRookCount = count
	return builder
}

func (builder *MaterialCountBuilder) WithBlackQueenCount(count uint8) *MaterialCountBuilder {
	builder.materialCount.BlackQueenCount = count
	return builder
}

func (builder *MaterialCountBuilder) Build() *MaterialCount {
	return builder.materialCount
}
