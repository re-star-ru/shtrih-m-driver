package driver

import (
	"errors"
	"fmt"
	"strings"
)

var (
	FSPrinterError01 = PrinterError{msg: "ФН: Неизвестная команда, неверный формат посылки или неизвестные параметры", num: 1} //FSPrinterError01 = ФН: Неизвестная команда, неверный формат посылки или неизвестные параметры
	FSPrinterError02 = PrinterError{msg: "ФН: Неверное состояние ФН", num: 2}                                                  //FSPrinterError02 = ФН: Неверное состояние ФН
	FSPrinterError03 = PrinterError{msg: "ФН: Ошибка ФН", num: 3}                                                              //FSPrinterError03 = ФН: Ошибка ФН
	FSPrinterError04 = PrinterError{msg: "ФН: Ошибка КС", num: 4}                                                              //FSPrinterError04 = ФН: Ошибка КС
	FSPrinterError05 = PrinterError{msg: "ФН: Закончен срок эксплуатации ФН", num: 5}                                          //FSPrinterError05 = ФН: Закончен срок эксплуатации ФН
	FSPrinterError06 = PrinterError{msg: "ФН: Архив ФН переполнен", num: 6}                                                    //FSPrinterError06 = ФН: Архив ФН переполнен

	PrinterError33 = PrinterError{msg: "Некорректные параметры в команде", num: 51} // PrinterError33 = Некорректные параметры в команде
	PrinterError40 = PrinterError{msg: "Переполнение диапазона скидок", num: 64}    // PrinterError40 = Переполнение диапазона скидок
	PrinterError4F = PrinterError{msg: "Неверный пароль", num: 79}
	PrinterError50 = PrinterError{msg: "Идет печать предыдущей команды", num: 80} // PrinterError50 = Идет печать предыдущей команды
	PrinterError69 = PrinterError{msg: "Переполнение денег по обороту налогов", num: 105}
	PrinterError73 = PrinterError{msg: "Команда не поддерживается в данном режиме", num: 115} // Команда не поддерживается в данном режиме
	PrinterError7E = PrinterError{msg: "Неверное значение в поле длины", num: 126}            //
	PrinterError8E = PrinterError{msg: "Нулевой итог чека", num: 142}                         // Нулевой итог чека

	PrinterError45 = PrinterError{msg: "Cумма всех типов оплаты меньше итога чека", num: 69}

	PrinterErrorUnknown = PrinterError{msg: "Неизвестная ошибка"}
)

var printerErrors = []PrinterError{
	PrinterError45,
	PrinterError69,
	PrinterError73,
	PrinterError8E,
	FSPrinterError01,
	FSPrinterError02,
	FSPrinterError03,
	FSPrinterError04,
	FSPrinterError05,
	FSPrinterError06,
	PrinterError40,
	PrinterError7E,
	PrinterError33,
	PrinterError50,
	PrinterError4F,
	PrinterErrorUnknown,
}

type PrinterError struct {
	msg string
	num byte
	err error
}

func (err PrinterError) Error() string {
	if err.err != nil {
		return fmt.Sprintf("num: %v, msg: %v, err: %v", err.num, err.msg, err.err)
	}

	return fmt.Sprintf("num: %v, msg: %v", err.num, err.msg)
}
func (err PrinterError) wrap(inner error) error {
	return PrinterError{msg: err.msg, err: inner}
}
func (err PrinterError) Unwrap() error {
	return err.err
}
func (err PrinterError) Is(target error) bool {
	ts := target.Error()
	return ts == err.msg || strings.HasPrefix(ts, err.msg+": ")
}

func checkOnPrinterError(err byte) error {
	if err == 0 {
		return nil
	}

	for _, v := range printerErrors {
		if err == v.num {
			return v
		}
	}

	return PrinterErrorUnknown.wrap(errors.New(fmt.Sprint("unknown error code:", err)))
}

//FSPrinterError08 = ФН: Нет запрошенных данных
//FSPrinterError09 = ФН: Некорректное значение параметров команды
//FSPrinterError10 = ФН: Превышение размеров TLV данных
//FSPrinterError11 = ФН: Нет транспортного соединения
//FSPrinterError12 = ФН: Исчерпан ресурс КС (криптографического сопроцессора)
//FSPrinterError14 = ФН: Исчерпан ресурс хранения
//FSPrinterError15 = ФН: Исчерпан ресурс ожидания передачи сообщения
//FSPrinterError16 = ФН: Продолжительность смены более 24 часов
//FSPrinterError17 = ФН: Неверная разница во времени между 2 операциями
//FSPrinterError20 = ФН: Сообщение от ОФД не может быть принято
//PrinterError00 = Ошибок нет
//PrinterError01 = ФП: Неисправен накопитель ФП 1, ФП 2 или часы
//PrinterError02 = ФП: Отсутствует ФП 1
//PrinterError03 = ФП: Отсутствует ФП 2
//PrinterError04 = ФП: Некорректные параметры в команде обращения к ФП
//PrinterError05 = ФП: Нет запрошенных данных
//PrinterError06 = ФП: ФП в режиме вывода данных
//PrinterError07 = ФП: Некорректные параметры в команде для данной реализации ФП
//PrinterError08 = ФП: Команда не поддерживается в данной реализации ФП
//PrinterError09 = ФП: Некорректная длина команды
//PrinterError0A = ФП: Формат данных не BCD
//PrinterError0B = ФП: Неисправна ячейка памяти ФП при записи итога
//PrinterError0C = ФП: Переполнение необнуляемой суммы
//PrinterError0D = ФП: Переполнение суммы итогов смен
//PrinterError11 = ФП: Не введена лицензия
//PrinterError12 = ФП: Заводской номер уже введен
//PrinterError13 = ФП: Текущая дата меньше даты последней записи в ФП
//PrinterError14 = ФП: Область сменных итогов ФП переполнена
//PrinterError15 = ФП: Смена уже открыта
//PrinterError16 = ФП: Смена не открыта
//PrinterError17 = ФП: Номер первой смены больше номера последней смены
//PrinterError18 = ФП: Дата первой смены больше даты последней смены
//PrinterError19 = ФП: Нет данных в ФП
//PrinterError1A = ФП: Область перерегистраций в ФП переполнена
//PrinterError1B = ФП: Заводской номер не введен
//PrinterError1C = ФП: В заданном диапазоне есть поврежденная запись
//PrinterError1D = ФП: Повреждена последняя запись сменных итогов
//PrinterError1E = ФП: Запись фискализации в накопителе не найдена
//PrinterError1F = ФП: Отсутствует память регистров
//PrinterError20 = ФП: Переполнение денежного регистра при добавлении
//PrinterError21 = ФП: Вычитаемая сумма больше содержимого денежного регистра
//PrinterError22 = ФП: Неверная дата
//PrinterError23 = ФП: Нет записи активизации
//PrinterError24 = ФП: Область активизаций переполнена
//PrinterError25 = ФП: Нет активизации с запрашиваемым номером
//PrinterError26 = ФП: В ФП больше 3 поврежденных записей
//PrinterError27 = ФП: Повреждение контрольных сумм ФП
//PrinterError28 = ФП: Технологическая метка в накопителе присутствует
//PrinterError29 = ФП: Технологическая метка в накопителе отсутствует
//PrinterError2A = ФП: Емкость микросхемы накопителя не соответствует текущей версии ПО
//PrinterError2B = Невозможно отменить предыдущую команду
//PrinterError2C = Обнулённая касса (повторное гашение невозможно)
//PrinterError2D = Сумма чека по секции меньше суммы сторно
//PrinterError2E = В ККТ нет денег для выплаты
//PrinterError2F = ФН: Таймаут связи
//PrinterError30 = ФН: ФН ответила NAK
//PrinterError31 = ФН: Ошибка в формате обмена
//PrinterError32 = ФН: Ошибка CRC
//PrinterError33 = Некорректные параметры в команде
//PrinterError34 = Нет данных
//PrinterError35 = Некорректный параметр при данных настройках
//PrinterError36 = Некорректные параметры в команде для данной реализации
//PrinterError37 = Команда не поддерживается в данной реализации
//PrinterError38 = Ошибка в ПЗУ
//PrinterError39 = Внутренняя ошибка ПО
//PrinterError3A = Переполнение накопления по надбавкам в смене
//PrinterError3B = Переполнение накопления в смене
//PrinterError3C = ЭКЛЗ: Неверный регистрационный номер
//PrinterError3D = Смена не открыта - операция невозможна
//PrinterError3E = Переполнение накопления по секциям в смене
//PrinterError3F = Переполнение накопления по скидкам в смене
//PrinterError40 = Переполнение диапазона скидок
//PrinterError41 = Переполнение диапазона оплаты наличными
//PrinterError42 = Переполнение диапазона оплаты типом 2
//PrinterError43 = Переполнение диапазона оплаты типом 3
//PrinterError44 = Переполнение диапазона оплаты типом 4
//PrinterError45 = Cумма всех типов оплаты меньше итога чека
//PrinterError46 = Не хватает наличности в кассе
//PrinterError47 = Переполнение накопления по налогам в смене
//PrinterError48 = Переполнение итога чека
//PrinterError49 = Операция невозможна в открытом чеке данного типа
//PrinterError4A = Открыт чек - операция невозможна
//PrinterError4B = Буфер чека переполнен
//PrinterError4C = Переполнение накопления по обороту налогов в смене
//PrinterError4D = Вносимая безналичной оплатой сумма больше суммы чека
//PrinterError4E = Смена превысила 24 часа
//PrinterError4F = Неверный пароль
//PrinterError50 = Идет печать предыдущей команды
//PrinterError51 = переполнение накоплений наличными в смене
//PrinterError52 = переполнение накоплений по типу оплаты 2 в смене
//PrinterError53 = переполнение накоплений по типу оплаты 3 в смене
//PrinterError54 = переполнение накоплений по типу оплаты 4 в смене
//PrinterError55 = Чек закрыт - операция невозможна
//PrinterError56 = Нет документа для повтора
//PrinterError57 = ЭКЛЗ: Количество закрытых смен не совпадает с ФП
//PrinterError58 = Ожидание команды продолжения печати
//PrinterError59 = Документ открыт другим оператором
//PrinterError5A = Скидка превышает накопления в чеке
//PrinterError5B = Переполнение диапазона надбавок
//PrinterError5C = Понижено напряжение 24В
//PrinterError5D = Таблица не определена
//PrinterError5E = Некорректная операция
//PrinterError5F = Отрицательный итог чека
//PrinterError60 = Переполнение при умножении
//PrinterError61 = Переполнение диапазона цены
//PrinterError62 = Переполнение диапазона количества
//PrinterError63 = Переполнение диапазона отдела
//PrinterError64 = ФП отсутствует
//PrinterError65 = Не хватает денег в секции
//PrinterError66 = Переполнение денег в секции
//PrinterError67 = Ошибка связи с ФП
//PrinterError68 = Не хватает денег по обороту налогов
//PrinterError69 = Переполнение денег по обороту налогов
//PrinterError6A = Ошибка питания в момент ответа по I2C
//PrinterError6B = Нет чековой ленты
//PrinterError6C = Нет контрольной ленты
//PrinterError6D = Не хватает денег по налогу
//PrinterError6E = Переполнение денег по налогу
//PrinterError6F = Переполнение по выплате в смене
//PrinterError70 = Переполнение ФП
//PrinterError71 = Ошибка отрезчика
//PrinterError72 = Команда не поддерживается в данном подрежиме
//PrinterError73 = Команда не поддерживается в данном режиме
//PrinterError74 = Ошибка ОЗУ
//PrinterError75 = Ошибка питания
//PrinterError76 = Ошибка принтера: нет импульсов с тахогенератора
//PrinterError77 = Ошибка принтера: нет сигнала с датчиков
//PrinterError78 = Замена ПО
//PrinterError79 = Замена ФП
//PrinterError7A = Поле не редактируется
//PrinterError7B = Ошибка оборудования
//PrinterError7C = Не совпадает дата
//PrinterError7D = Неверный формат даты
//PrinterError7E = Неверное значение в поле длины
//PrinterError7F = Переполнение диапазона итога
//PrinterError80 = Ошибка связи с ФП
//PrinterError81 = Ошибка связи с ФП
//PrinterError82 = Ошибка связи с ФП
//PrinterError83 = Ошибка связи с ФП
//PrinterError84 = Переполнение наличности
//PrinterError85 = Переполнение по продажам в смене
//PrinterError86 = Переполнение по покупкам в смене
//PrinterError87 = Переполнение по возвратам продаж в смене
//PrinterError88 = Переполнение по возвратам покупок в смене
//PrinterError89 = Переполнение по внесению в смене
//PrinterError8A = Переполнение по надбавкам в чеке
//PrinterError8B = Переполнение по скидкам в чеке
//PrinterError8C = Отрицательный итог надбавки в чеке
//PrinterError8D = Отрицательный итог скидки в чеке
//PrinterError8E = Нулевой итог чека
//PrinterError8F = Касса не фискализирована
//PrinterError90 = Поле превышает размер установленный в настройках
//PrinterError91 = Выход за границу поля печати при данных настройках шрифта
//PrinterError92 = Наложение полей
//PrinterError93 = Восстановление ОЗУ прошло успешно
//PrinterError94 = Исчерпан лимит операций в чеке
//PrinterError95 = Неизвестная ошибка ЭКЛЗ
//PrinterError96 = Выполните суточный отчет с гашением
//PrinterError9B = Некорректное действие
//PrinterError9C = Товар не найден по коду в базе товаров
//PrinterError9D = Неверные данные в записе о товаре в базе товаров
//PrinterError9E = Неверный размер файла базы или регистров товаров
//PrinterErrorA0 = Ошибка связи с ЭКЛЗ
//PrinterErrorA1 = ЭКЛЗ отсутствует
//PrinterErrorA2 = ЭКЛЗ: Некорректный формат или параметр команды
//PrinterErrorA3 = ЭКЛЗ: Некорректное состояние ЭКЛЗ
//PrinterErrorA4 = ЭКЛЗ: Авария ЭКЛЗ
//PrinterErrorA5 = ЭКЛЗ: Авария КС в составе ЭКЛЗ
//PrinterErrorA6 = ЭКЛЗ: Исчерпан временной ресурс ЭКЛЗ
//PrinterErrorA7 = ЭКЛЗ: ЭКЛЗ переполнена
//PrinterErrorA8 = ЭКЛЗ: Неверные дата или время
//PrinterErrorA9 = ЭКЛЗ: Нет запрошенных данных
//PrinterErrorAA = ЭКЛЗ: Переполнение ЭКЛЗ (отрицательный итог документа)
//PrinterErrorAF = Некорректные значения принятых данных от ЭКЛЗ
//PrinterErrorB0 = ЭКЛЗ: Переполнение в параметре количество
//PrinterErrorB1 = ЭКЛЗ: Переполнение в параметре сумма
//PrinterErrorB2 = ЭКЛЗ: Уже активизирована
//PrinterErrorB4 = Найденная запись фискализации повреждена
//PrinterErrorB5 = Запись заводского номера ККМ повреждена
//PrinterErrorB6 = Найденная запись активизации ЭКЛЗ повреждена
//PrinterErrorB7 = Записи сменных итогов в накопителе не найдены
//PrinterErrorB8 = Последняя запись сменных итогов не записана
//PrinterErrorB9 = Сигнатура версии структуры данных в накопителе не совпадает с текущей версией ПО
//PrinterErrorBA = Структура накопителя повреждена
//PrinterErrorBB = Текущая дата меньше даты последней записи активизации ЭКЛЗ
//PrinterErrorBC = Текущая дата меньше даты последней записи фискализации
//PrinterErrorBD = Текущая дата меньше даты последней записи сменного итога
//PrinterErrorBE = Команда не поддерживается в текущем состоянии
//PrinterErrorBF = Инициализация накопителя невозможна
//PrinterErrorC0 = Контроль даты и времени (подтвердите дату и время)
//PrinterErrorC1 = ЭКЛЗ: суточный отчет с гашением прервать нельзя
//PrinterErrorC2 = Превышение напряжения блока питания
//PrinterErrorC3 = Несовпадение итогов чека с ЭКЛЗ
//PrinterErrorC4 = Несовпадение номеров смен
//PrinterErrorC5 = Буфер подкладного документа пуст
//PrinterErrorC6 = Подкладной документ отсутствует
//PrinterErrorC7 = Поле не редактируется в данном режиме
//PrinterErrorC8 = Ошибка связи с принтером
//PrinterErrorC9 = Перегрев печатающей головки
//PrinterErrorCA = Температура вне условий эксплуатации
//PrinterErrorCB = Неверный подытог чека
//PrinterErrorCC = Смена в ЭКЛЗ уже закрыта
//PrinterErrorCD = Тест целостности архива ЭКЛЗ не прошел
//PrinterErrorCE = Объем ОЗУ или ПЗУ на ККМ исчерпан
//PrinterErrorCF = Неверная дата (Часы сброшены? Установите дату!)
//PrinterErrorD0 = Не распечатана контрольная лента по смене из ЭКЛЗ
//PrinterErrorD1 = Нет документов в буфере
//PrinterErrorD2 = Модем не работает
//PrinterErrorD3 = Товар не произведен или выбыл
//PrinterErrorD4 = Код маркировки фальсифицирован
//PrinterErrorD5 = Ошибка аутентификации
//PrinterErrorE0 = Ошибка связи с купюроприемником
//PrinterErrorE1 = Купюроприемник занят
//PrinterErrorE2 = Итог чека не соответствует итогу купюроприемника
//PrinterErrorE3 = Ошибка купюроприемника
//PrinterErrorE4 = Итог купюроприемника не нулевой
//PrinterError0100 = Открыта крышка
//PrinterError0101 = Нет бумаги или открыта крышка
