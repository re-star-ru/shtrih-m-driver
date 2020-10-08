package shtrih

import "shtrih-drv/internal/shtrih/FDTags"

func (p *Printer) PrintCheck() {
	p.Ping()

	// todo print check
}

// 1
// Открытие смены, передаются данные о кассире
func (p *Printer) OpenShift(cashier string) {
	p.Ping()

	// todo print check
}

// last
// закрытие смены, так же передаются данные о кассире
func (p *Printer) CloseShift(cashier string) {
	p.Ping()

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

	// 1021, кассир
	//p.writeCachierName(params.Parameters.CachierName);  // печать имени кассира, возможно инн?

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

	writeVATINTagIfNotNullAndNotEmpty(p, FDTags.CashierINN, params.Parameters.CashierVATIN)

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

//
//prepare(printer);
//
//LongPrinterStatus printerStatus = printer.readLongPrinterStatus();
//
//// данная проверка нужна, т.к. код начала чека автоматом открывает смену, при этом
//// в открытие смены не будет передан ИНН кассира
//if (printerStatus.getPrinterMode().isDayClosed()) {
//throw JposExceptionHandler.getJposException(
//new SmFiscalPrinterException(2, "Закрытая смена, операция невозможна"));
//}
//
//if (isCashCore(printer)) {
//printer.writeTable(1, 1, 48, electronically ? "1" : "0");
//}
//
//// 1021, кассир
//printer.writeCashierName(params.Parameters.CashierName);
//
//// Задаем тип чека
//if (params.Parameters.PaymentType == 1) { // приход
//printer.setFiscalReceiptType(SmFptrConst.SMFPTR_RT_SALE);
//} else if (params.Parameters.PaymentType == 2) { // возврат прихода
//printer.setFiscalReceiptType(SmFptrConst.SMFPTR_RT_RETSALE);
//} else if (params.Parameters.PaymentType == 3) { // расход
//printer.setFiscalReceiptType(SmFptrConst.SMFPTR_RT_BUY);
//} else if (params.Parameters.PaymentType == 4) { // возврат расхода
//printer.setFiscalReceiptType(SmFptrConst.SMFPTR_RT_RETBUY);
//} else {
//throw new UnsupportedOperationException("Неизвестный тип чека " + params.Parameters.PaymentType);
//}
//
//// Указываем систему налогообложения
//printer.setParameter(SmFptrConst.SMFPTR_DIO_PARAM_TAX_SYSTEM, params.Parameters.getTaxVariant());
//
//printer.beginFiscalReceipt(false);
//
//FSStatusInfo fsStatus = printer.fsReadStatus();
//int docNumebr = (int) fsStatus.getDocNumber() + 1; // Номер ФД
//int shiftNumber = printer.readLongPrinterStatus().getCurrentShiftNumber(); // Номер смены
//
//if (electronically) {
//if (isDesktop(printer) || isShtrihNano(printer)) {
//printer.writeTable(17, 1, 7, "1");
//}
//}
//
//if (isCashCore(printer)) {
//printer.writeTable(1, 1, 48, electronically ? "1" : "0");
//}
//
//boolean isNewMobile = isMobile(printer) && printerStatus.getFirmwareBuild() >= 20041;
//if (isNewMobile) {
//printer.writeTable(1, 1, 43, electronically ? "1" : "0");
//}
//
//boolean sendOperationTagsFirst = isCashCore(printer) || isMobile(printer);
//
//for (Object element : params.Positions.FiscalStrings) {
//
//if (element instanceof FiscalString) {
//
//FiscalString item = (FiscalString) element;
//
//if(!sendOperationTagsFirst)
//fsOperationV2(printer, item);
//
//// Единица измерения предмета расчета, 1197
//if (item.MeasurementUnit != null) {
//printer.fsWriteOperationTag(1197, item.MeasurementUnit);
//}
//
//// Признак агента по предмету расчета, 1222
//if (item.SignSubjectCalculationAgent > 0) {
//printer.fsWriteOperationTag(1222, item.SignSubjectCalculationAgent, 1);
//}
//
//// Данные агента, 1223
//if (item.AgentData != null) {
//
//byte[] agentData = buildAgentDataTLV(item.AgentData);
//
//if (agentData.length > 0)
//printer.fsWriteOperationTag(1223, agentData);
//}
//
//// Информация о поставщике
//if (item.PurveyorData != null) {
//
//// ИНН поставщика, 1226
//if (item.PurveyorData.PurveyorVATIN != null)
//printer.fsWriteOperationTag(1226, item.PurveyorData.PurveyorVATIN);
//
//byte[] purveyorData = buildPurveyorDataTLV(item.PurveyorData);
//
//// Данные поставщика, 1224
//if (purveyorData.length > 0)
//printer.fsWriteOperationTag(1224, purveyorData);
//}
//
//// Код товарной номенклатуры, 1162
//if (item.GoodCodeData != null) {
//byte[] data = generateKTN(item.GoodCodeData);
//printer.fsWriteOperationTag(1162, data);
//}
//
//// TODO: Дополнительный реквизит предмета расчета, 1191
//// TODO: Код страны происхождения товара, 1230
//// TODO: Номер таможенной декларации, 1231
//// TODO: Акциз, 1229
//
//// В мобайле и КЯ сперва нужно отправить тэги затем отправлять Операцию V2
//if (sendOperationTagsFirst)
//fsOperationV2(printer, item);
//
//} else if (element instanceof TextString) {
//TextString item = (TextString) element;
//printer.printRecMessage(item.Text, item.FontNumber);
//} else if (element instanceof Barcode) {
//Barcode item = (Barcode) element;
//printBarCode(printer, item);
//} else {
//throw new UnsupportedOperationException("Unknown element of type " + element.getClass().getName());
//}
//}
//
//// ИНН кассира
//writeVATINTagIfNotNullAndNotEmpty(printer, 1203, params.Parameters.CashierVATIN);
//
//// телефон или электронный адрес покупателя, не могут быть одновременно заданы
//writeTagIfNotNullAndNotEmpty(printer, 1008, params.Parameters.CustomerEmail);
//writeTagIfNotNullAndNotEmpty(printer, 1008, params.Parameters.CustomerPhone);
//
//// адрес электронной почты отправителя чека
//writeTagIfNotNullAndNotEmpty(printer, 1117, params.Parameters.SenderEmail);
//
//// адрес расчетов, пердача данного тэга возможна только в режиме развозной торговли
//writeTagIfNotNullAndNotEmpty(printer, 1009, params.Parameters.AddressSettle);
//
//// место расчетов, пердача данного тэга возможна только в режиме развозной торговли
//writeTagIfNotNullAndNotEmpty(printer, 1187, params.Parameters.PlaceSettle);
//
//// признак агента
//if (params.Parameters.AgentSign > 0)
//printer.fsWriteTag(1057, params.Parameters.AgentSign, 1);
//
//// операция платежного агента
//writeTagIfNotNullAndNotEmpty(printer, 1044, params.Parameters.AgentData.PayingAgentOperation);
//// телефон платежного агента
//writeTagIfNotNullAndNotEmpty(printer, 1073, params.Parameters.AgentData.PayingAgentPhone);
//// телефон оператора по приему платежей
//writeTagIfNotNullAndNotEmpty(printer, 1074, params.Parameters.AgentData.ReceivePaymentsOperatorPhone);
//// телефон оператора перевода
//writeTagIfNotNullAndNotEmpty(printer, 1075, params.Parameters.AgentData.MoneyTransferOperatorPhone);
//// наименование оператора перевода
//writeTagIfNotNullAndNotEmpty(printer, 1026, params.Parameters.AgentData.MoneyTransferOperatorName);
//// адрес оператора перевода
//writeTagIfNotNullAndNotEmpty(printer, 1005, params.Parameters.AgentData.MoneyTransferOperatorAddress);
//// ИНН оператора перевода
//writeVATINTagIfNotNullAndNotEmpty(printer, 1016, params.Parameters.AgentData.MoneyTransferOperatorVATIN);
//
//// телефон поставщика
//writeTagIfNotNullAndNotEmpty(printer, 1171, params.Parameters.PurveyorData.PurveyorPhone);
//
//// TODO: Номер автомата, 1036
//// TODO: Дополнительный реквизит пользователя, 1084
//// TODO: Дополнительный реквизит чека (БСО), 1192
//// TODO: Покупатель (клиент), 1227
//// TODO: ИНН покупателя (клиента), 1228
//
//if (params.Payments.getCash() > 0)
//printer.printRecTotal(0, params.Payments.getCash(), "0");
//if (params.Payments.getElectronicPayment() > 0)
//printer.printRecTotal(0, params.Payments.getElectronicPayment(), "1");
//if (params.Payments.getAdvancePayment() > 0)
//printer.printRecTotal(0, params.Payments.getAdvancePayment(), "13");
//if (params.Payments.getCredit() > 0)
//printer.printRecTotal(0, params.Payments.getCredit(), "14");
//if (params.Payments.getCashProvision() > 0)
//printer.printRecTotal(0, params.Payments.getCashProvision(), "15");
//
//printer.endFiscalReceipt(false);
//}
