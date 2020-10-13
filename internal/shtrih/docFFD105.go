package shtrih

import (
	"errors"
	"fmt"
	"shtrih-drv/internal/shtrih/check"
	"shtrih-drv/internal/shtrih/tables"
)

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

	chk := check.New()

	// 1021, кассир
	//p.writeCashierName(params.CashierName) // печать имени кассира, возможно инн?
	//p.WriteTable(tables.SmfpTableCashier, 14, 2, "Оператор14") // запись имени текущего кассира в таблицу
	chk.CashierName = "Оператор 14"

	// Задаем тип чека
	switch params.PaymentType {
	case 1: // приход
		chk.FiscalReceiptType = check.SMFPTR_RT_SALE
	case 2: // возврат прихода
		chk.FiscalReceiptType = check.SMFPTR_RT_RETSALE
	case 3: // расход
		chk.FiscalReceiptType = check.SMFPTR_RT_BUY
	case 4: // возврат расхода
		chk.FiscalReceiptType = check.SMFPTR_RT_RETBUY
	default:
		err := errors.New(fmt.Sprint("Неизвестный тип чека: ", params.PaymentType))
		p.logger.Fatal(err)
		return err
	}

	// Указываем систему налогообложения
	//p.SetParameter(SMFPTR_DIO_PARAM_TAX_SYSTEM = 16, params.Parameters.getTaxVariant())
	chk.TaxVariant = params.TaxVariant

	//p.beginFiscalReceipt(false) // начинаем печать чека без хедера

	//fsStatus := p.FnReadStatus() // читаем статус кассы
	//docNumber := fsStatus.getDocumentNumber + 1                            // Получаем номер создаваемого фискального документа
	//shiftNumber := printer.readLongPrinterStatus().getCurrentShiftNumber() // получаем номер текущей смены

	if electronically {
		p.WriteTable(tables.RegionalSettings, 1, 7, "1") // не печатать 1 документ
	}

	//52011 номер сборки
	//boolean isNewMobile = isMobile(printer) && printerStatus.getFirmwareBuild() >= 20041;
	//if (isNewMobile) {
	//	printer.writeTable(1, 1, 43, electronically ? "1" : "0"); // авто печать журнала? только для мобильных принтеров
	//}

	//boolean sendOperationTagsFirst = isCashCore(printer) || isMobile(printer); отправлять первым тег операции, только для мобильных и там где есть кассовое ядро(?)

	//err := p.writeFiscalStrings(params.Positions.FiscalStrings)

	//Инн кассира
	//writeVATINTagIfNotNullAndNotEmpty(p, FDTags.CashierINN, params.Parameters.CashierVATIN)
	p.FNWriteTLV()

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

//КАССОВЫЙ ЧЕК
//ФН:9281000100007442
//РН ККТ:0001837854048714
//ИНН:263209745357
//ФД:11904
//ДАТА, ВРЕМЯ:13.10.2020 15:34:00
//ФП:4074622255 (3104F2DDCD2F)
//СМЕНА:920
//НОМЕР ЧЕКА ЗА СМЕНУ:6
//ПРИЗН. РАСЧЕТА:1 (Приход)
//ИТОГ:1500.00
//ИНН КАССИРА:263209745357
//ПРЕДМ. РАСЧЕТА:
//НАИМЕН. ПРЕДМ. РАСЧЕТА:Дефектовка
//ЦЕНА ЗА ЕД. ПРЕДМ. РАСЧ.:1500.00
//КОЛ-ВО ПРЕДМ. РАСЧЕТА:1.000000
//СТОИМ. ПРЕДМ. РАСЧЕТА:1500.00
//ПРИЗН. СПОСОБА РАСЧ.:4
//ПРИЗН. ПРЕДМЕТА РАСЧ.:4
//КАССИР:Волков Е. И.
//НАЛИЧНЫМИ:1500.00
//БЕЗНАЛИЧНЫМИ:0.00
//ПРЕДВАРИТЕЛЬНАЯ ОПЛАТА (АВАНС):0.00
//ПОСЛЕДУЮЩАЯ ОПЛАТА (КРЕДИТ):0.00
//ИНАЯ ФОРМА ОПЛАТЫ:0.00
//САЙТ ФНС:www.nalog.ru
//МЕСТО РАСЧЕТОВ:касса
//ВЕРСИЯ ФФД:2 (1.05)
//СУММА БЕЗ НДС:1500.00
//НАИМЕН. ПОЛЬЗ.:ИП Волков Евгений Иванович
//АДР.РАСЧЕТОВ:Ставропольский край, пос. Горячеводский, ул. Совхозная д. 85
//СНО:8 (ЕНВД)
