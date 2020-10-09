package shtrih

func (p *Printer) PrintCheck() {
	//p.Ping()

	// todo print check
}

// 1
// Открытие смены, передаются данные о кассире
func (p *Printer) OpenShift(cashier string) {
	//p.Ping()

	// todo print check
}

// last
// закрытие смены, так же передаются данные о кассире
func (p *Printer) CloseShift(cashier string) {
	//p.Ping()

	// todo print check
}

/**
 * Формирование чека
 *
 * @param params         данные чека
 * @param electronically Формирование чека в только электроном виде. Печать чека не осуществляется.
 */
func (p *Printer) ProcessCheck(params CheckPackage, electronically bool) error {

	// TODO: проверка на открытую смену
	//LongPrinterStatus printerStatus = printer.readLongPrinterStatus();
	//
	//// данная проверка нужна, т.к. код начала чека автоматом открывает смену, при этом
	//// в открытие смены не будет передан ИНН кассира
	//if (printerStatus.getPrinterMode().isDayClosed()) {
	//throw JposExceptionHandler.getJposException(
	//new SmFiscalPrinterException(2, "Закрытая смена, операция невозможна"));
	//}

	// 1021, кассир
	p.writeCashierName(params.CashierName) // печать имени кассира, возможно инн?

	// Задаем тип чека
	//// Receipt type constants
	//public static final int SMFPTR_RT_SALE = 100;
	//public static final int SMFPTR_RT_BUY = 101;
	//public static final int SMFPTR_RT_RETSALE = 102;
	//public static final int SMFPTR_RT_RETBUY = 103;
	//switch params.Parameters.PaymentType {
	//case 1: // приход
	//	printer.setFiscalReceiptType(SMFPTR_RT_SALE = 100)
	//case 2: // возврат прихода
	//	printer.setFiscalReceiptType(SMFPTR_RT_RETSALE = 102);
	//case 3: // расход
	//	printer.setFiscalReceiptType(SMFPTR_RT_BUY = 101);
	//case 4: // возврат расхода
	//	printer.setFiscalReceiptType(SMFPTR_RT_RETBUY = 103);
	//default:
	//	return errors.New("Неизвестный тип чека: ", params.Parameters.PaymentType)
	//}

	// Указываем систему налогообложения
	//p.SetParameter(SMFPTR_DIO_PARAM_TAX_SYSTEM = 16, params.Parameters.getTaxVariant())

	//p.beginFiscalReceipt(false) // начинаем печать чека без хедера

	//fsStatus := p.FnReadStatus() // читаем статус кассы
	//docNumber := fsStatus.getDocumentNumber + 1                            // Получаем номер создаваемого фискального документа
	//shiftNumber := printer.readLongPrinterStatus().getCurrentShiftNumber() // получаем номер текущей смены

	//if (electronically) {
	//		printer.writeTable(17, 1, 7, "1"); // не печатать 1 документ
	//}

	//52011 номер сборки
	//boolean isNewMobile = isMobile(printer) && printerStatus.getFirmwareBuild() >= 20041;
	//if (isNewMobile) {
	//	printer.writeTable(1, 1, 43, electronically ? "1" : "0"); // авто печать журнала? только для мобильных принтеров
	//}

	//boolean sendOperationTagsFirst = isCashCore(printer) || isMobile(printer); отправлять первым тег операции, только для мобильных и там где есть кассовое ядро(?)

	//err := p.writeFiscalStrings(params.Positions.FiscalStrings)

	//Инн кассира
	//writeVATINTagIfNotNullAndNotEmpty(p, FDTags.CashierINN, params.Parameters.CashierVATIN)

	// телефон или электронный адрес покупателя, не могут быть одновременно заданы
	//writeTagIfNotNullAndNotEmpty(printer, FDTags.ClientEmailOrNumber, params.Parameters.CustomerEmail)
	//writeTagIfNotNullAndNotEmpty(printer, FDTags.ClientEmailOrNumber, params.Parameters.CustomerPhone)

	// адрес электронной почты отправителя чека
	//writeTagIfNotNullAndNotEmpty(printer, FDTags.EmailCheckSender, params.Parameters.SenderEmail)

	// признак агента
	//if (params.Parameters.AgentSign > 0) {
	//	printer.fsWriteTag(1057, params.Parameters.AgentSign, 1);
	//}

	//if params.Payments.getCash() > 0 {
	//	printer.printRecTotal(0, params.Payments.getCash(), "0")
	//}
	//
	//if params.Payments.getElectronicPayment() > 0 {
	//	printer.printRecTotal(0, params.Payments.getElectronicPayment(), "1")
	//
	//}
	//
	//if params.Payments.getAdvancePayment() > 0 {
	//	printer.printRecTotal(0, params.Payments.getAdvancePayment(), "13")
	//
	//}
	//if params.Payments.getCredit() > 0 {
	//	printer.printRecTotal(0, params.Payments.getCredit(), "14")
	//
	//}
	//if params.Payments.getCashProvision() > 0 {
	//	printer.printRecTotal(0, params.Payments.getCashProvision(), "15")
	//}

	//p.endFiscalReceipt(false);
	return nil
}

func (p *Printer) writeFiscalStrings(fiscalStrings []string) error {
	return nil
}

func writeVATINTagIfNotNullAndNotEmpty(printer Printer, tagID int, value string) error {
	//if value != "" {
	//	printer.fsWriteTag(tagId, formatInn(value)) // write string tag
	//}
	return nil
}

func writeTagIfNotNullAndNotEmpty(printer Printer, tagId int, value string) error {
	//if value != "" {
	//	printer.fsWriteTag(tagId, value)
	//}

	return nil
}
