package resolve

var (
	NoOpDataTypeResolver       DataTypeResolver       = ChainDataTypeResolvers()
	NoOpSecretTypeResolver     SecretTypeResolver     = ChainSecretTypeResolvers()
	NoOpConnectionTypeResolver ConnectionTypeResolver = ChainConnectionTypeResolvers()
	NoOpOutputTypeResolver     OutputTypeResolver     = ChainOutputTypeResolvers()
	NoOpParameterTypeResolver  ParameterTypeResolver  = ChainParameterTypeResolvers()
	NoOpAnswerTypeResolver     AnswerTypeResolver     = ChainAnswerTypeResolvers()
)
