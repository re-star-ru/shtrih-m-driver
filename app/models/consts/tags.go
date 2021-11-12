package consts

type TLVTag struct {
	Name   string // наименование тега
	Type   string // тип тега
	Code   uint16 // код тега
	Length uint16 // длинна тега
}

var CashierINN = TLVTag{
	Name:   "ИНН Кассира",
	Type:   "String",
	Code:   1203,
	Length: 12,
} // Инн кассира

const (
	DocumentName        = 1000 // наименование документа
	VersionFFD          = 1209 // номер версии ФФД
	UserName            = 1048 // наименование пользователя
	UserINN             = 1018 // ИНН пользователя
	TaxationSystems     = 1062 // Системы налогообложения
	DateTime            = 1012 // Дата Время
	KKTRegisterNumber   = 1037 // регистрационный номер ККТ
	OfflineSign         = 1002 // признак автономного режима
	PrinterSetupSign    = 1221 // признак установки принтера в автомате
	ASBSO               = 1110 // Признак Автоматизированная Система для Бланков Строгой Отчетности
	EncryptionSign      = 1056 // Признак шифрования
	AutomationSign      = 1001 // Признак автоматического режима ?
	OnlyInternetKKTSign = 1108 // Признак ККТ только для интернет рассчетов
	NumberAutomat       = 1036 // Номер автомата мб сериыйнй номер ккт
	SellExcisableGoods  = 1207 // признак торговли подакцизными товарами
	KeyResourceFP       = 1213 // ресурс ключей ФП
	PaymentForService   = 1109 // признак расчетов за услуги
	AgentSign           = 1057 // признак агента
	KKTSerialNumber     = 1013 // заводской номер ККТ
	CashierName         = 1021 // Кассир (имя?)
	BillingAddress      = 1009 // Рассчетный адрес
	BillingSpace        = 1187 // Место рассчетов
	FNSSiteAdress       = 1060 // Адрес ФНС
	EmailCheckSender    = 1117 // email отправителя чека
	INNOFD              = 1017 // ИНН ОФД
	OFDName             = 1046 // Наименование офд
	KKTVersion          = 1188 // Версия Контрольно кассовой техники
	FFDKKTVersion       = 1189 // Версия формат фискального документа ККТ
	FFDFN               = 1190 // Версия формата фискального документа в фискальном накопителе
	FDNumber            = 1040 // Номер Фискального документа
	FNNumber            = 1041 // Номер фискального накопителя
	FPD                 = 1077 // Фискальный признак документа

	Client              = 1227 // Покупатель клиент
	ClientINN           = 1228 // ИНН покупателя
	CheckNumberForShift = 1042 // Номер чека за смену
	SettlementAttribute = 1054 // ПризнакРассчета
	TaxationSystemCheck = 1055 // Применяемая система налогообложения
	ClientEmailOrNumber = 1008 // емейл или номер клиента

	SubjectOfCalculation  = 1059 // предмет рассчета
	CalculationSum        = 1020 // сумма рассчета указанного в чеке
	CalculationSumCash    = 1031 // сумма по чеку наличными
	CalculationSumNonCash = 1081 // сумма по чеку безналичными

	CalculationSumWithVAT0   = 1104 // сумма расчета по чеку с НДС по ставке 0%
	CalculationSumWithoutVAT = 1105 // сумма расчета по чеку без НДС
	CalculationVAT20120      = 1106 // сумма ндс по растчетной ставке 20/120
	QRCode                   = 1196 // QR-код, обязательлный
)
