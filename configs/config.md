Параметры драйвера
Параметры драйвера хранятся в файле jpos.xml. Загрузкой параметров занимается ControlObject, то есть библиотека jpos. Драйвер получает JposEntry.
1.	Тип порта: 0 - последовательный порт, 1 - bluetooth, 2 - socket, 3 - создание класса по названию
<!-- Port type: 0 - serial, 1 - bluetooth, 2 - socket, 3 - from parameter protClass -->
<prop name="portType" type="String" value="0"/>
2.	Название класса порта
<!-- portClass  -->
<prop name="portClass" type="String" value="com.shtrih.fiscalprinter.port.SerialPrinterPort"/>
3.	Тип протокола: 0 - протокол 1.0, 1 - протокол версии 2.0
<!-- ProtocolType, 0 - protocol 1, 1 - protocol 2 -->
<prop name="protocolType" type="String" value="0"/>
4.	Имя порта.
<!--Port name-->
<prop name="portName" type="String" value="COM1"/>
5.	Скорость связи
<prop name="baudRate" type="String" value="115200"/>
6.	Отдел по умолчанию. Можно изменить также через directIO
<!--Default department-->
<prop name="department" type="String" value="1"/>
7.	Номер шрифта по умолчанию. Можно изменить также через directIO
<!-- Default font number -->
<prop name="fontNumber" type="String" value="1"/>
8.	Текст для закрытия чека
<!-- Close receipt text -->
<prop name="closeReceiptText" type="String" value=""/>
9.	Текст подитога, печатается в методе printRecSubtotal
<!-- Subtotal text -->
<prop name="subtotalText" type="String" value="ПОДИТОГ:"/>
10.	Таймаут приема байта драйвера в миллисекундах.
<!-- Driver byte receive timeout -->
<prop name="byteTimeout" type="String" value="3000"/>
11.	Таймаут приема байта устройства в миллисекундах, записывается в ФР при инициализации
<!-- Device byte receive timeout -->
<prop name="deviceByteTimeout" type="String" value="3000"/>
12.	Разрешение поиска устройства на всех портах системы. Драйвер начинает поиск устройства, если не удалось подключиться к устройству с заданными параметрами. По умолчанию выключен. 
<!-- Device search enabled for all serial ports -->
<prop name="searchByPortEnabled" type="String" value="0"/>
13.	Разрешение поиска устройства на всех скоростях. По умолчанию включен.
<!-- Device search enabled for all baud rates -->
<prop name="searchByBaudRateEnabled" type="String" value="1"/>
14.	Пароль налогового инспектора
<!-- Tax officer password -->
<prop name="taxPassword" type="String" value="0"/>
15.	Пароль оператора
<!-- Operator password -->
<prop name="operatorPassword" type="String" value="1"/>
16.	Пароль системного администратора
<!-- System administrator password -->
<prop name="sysAdminPassword" type="String" value="30"/>
17.	Разрешение опроса устройства. Опрос устройства нужен для оповещения приложения о 
<!-- device state polling enabled -->
<prop name="pollEnabled" type="String" value="0"/>
18.	Интервал опроса в миллисекундах
<!-- device state polling interval in milliseconds -->
<prop name="pollInterval" type="String" value="100"/>
19.	Коэффициент для сумм
<!-- Amount coefficient -->
<prop name="amountFactor" type="String" value="1"/>
20.	Коэффициент для количества
<!-- Quantity coefficient -->
<prop name="quantityFactor" type="String" value="1"/>
21.	Кодировка текстовых строк
<!-- Strings encoding -->
<prop name="stringEncoding" type="String" value="Cp866"/>
22.	Имя файла статистики
<!-- Statistics file name -->
<prop name="statisticFileName" type="String" value="ShtrihFiscalPrinter.xml"/>
23.	Задержка после печати штрихкода командой "Печать графической линии"
<!-- Barcode print time -->
<prop name="graphicsLineDelay" type="String" value="1000"/>
24.	Имя файла для записи таблиц ФР. Драйвер запишет эти данные в ФР при инициализации
<!-- fieldsFileName to initialize fiscalprinter tables -->
<prop name="fieldsFileName" type="String" value="tables.csv"/>
25.	Имя папки для поиска файлов таблиц. Драйвер выбирает файл по имени устройства. 
<!-- fieldsFilesPath to initialize fiscalprinter tables -->
<prop name="fieldsFilesPath" type="String" value="I:\Projects\JavaPOS\Bin\tables\"/>
26.	Количество строк заголовка чека
<!-- Number of header lines -->
<prop name="numHeaderLines" type="String" value="4"/>
27.	Количество строк рекламного текста
<!-- Number of trailer lines -->
<prop name="numTrailerLines" type="String" value="3"/>
28.	Тип отчета, который снимается при вызове printReport: 0 - отчет по ФП, 1 - отчет по ЭКЛЗ
<!-- Device to print report, 0 - fiscal memory (FM), 1 - electronic journal (EJ)  -->
<prop name="reportDevice" type="String" value="0"/>
29.	Тип отчета, 0 - краткий, 1 - полный
<!-- Report type, 0 - short, 1 - full -->
<prop name="reportType" type="String" value="1"/>
30.	Команда для чтения состояния ФР: 0 - краткий запрос состояния 10h, 1 - полный запрос, 2 - автоматический выбор (10h, если поддерживается, иначе 11h)
<!-- Status command:   -->
<!-- 0 - command 10h, short status request  -->
<!-- 1 - command 11h, long status request  -->
<!-- 2 - status command selected by driver -->
<prop name="statusCommand" type="String" value="0"/>
31.	Файл для всех сообщений драйвера
<!-- Localization file name -->
<prop name="messagesFileName" type="String" value="shtrihjavapos_en.properties"/>
32.	Разрешение переноса длинных строк
<!-- Wrap text enabled -->
<prop name="wrapText" type="String" value="1"/>
33.	Задержка после закрытия чека. Может потребоваться для некоторых моделей ФР. 
<!-- Sleep time after receipt close -->
<prop name="recCloseSleepTime" type="String" value="0"/>
34.	Тип отрезки чека, 0 - полная, 1 - неполная
<!-- Cut type, 0 - full cut, 1 - partial cut -->
<prop name="cutType" type="String" value="1"/>
35.	Режим отрезки, 0 - автоматически, 1 - отрезка запрещена
<!-- Cut mode, 0 - auto, 1 - disabled -->
<prop name="cutMode" type="String" value="0"/>
36.	Параметр протокола ФР 1.0, максимальное количество запросов ENQ при передаче одной команды
<!-- maxEnqNumber -->
<prop name="maxEnqNumber" type="String" value="3"/>
37.	Параметр протокола ФР 1.0, максимальное количество ответов NAK при передаче одной команды, то есть максимальное количество ошибок при передаче команды.
<!-- maxNakCommandNumber -->
<prop name="maxNakCommandNumber" type="String" value="3"/>
38.	Параметр протокола ФР 1.0, максимальное количество ответов NAK при приеме ответа.
<!-- maxNakAnswerNumber -->
<prop name="maxNakAnswerNumber" type="String" value="3"/>
39.	Максимальное количество повторов команды
<!-- maxRepeatCount -->
<prop name="maxRepeatCount" type="String" value="1"/>
40.	Названия типов оплаты приложения
<!-- Payment types -->
<prop name="payType0" type="String" value="0"/>
<prop name="payType39" type="String" value="3"/>
41.	Названия типов оплаты ФР
<!-- Payment names -->
<prop name="paymentName1" type="String" value="CASH"/>
<prop name="paymentName2" type="String" value="CREDIT"/>
<prop name="paymentName3" type="String" value="КАРТА"/>
<prop name="paymentName4" type="String" value="СКИДКА"/>
42.	Названия налоговых групп
<!-- Tax names -->
<prop name="taxName0" type="String" value="НДС 10%"/>
<prop name="taxName1" type="String" value="НДС 18%"/>
<prop name="taxName2" type="String" value="C"/>
<prop name="taxName3" type="String" value="D"/>
43.	Разрешение получения Z- отчета в виде XML файла. Драйвер сохраняет данные всех денежных и операционных регистров XML файле. 
<!-- create Z-report in XML format -->
<prop name="XmlZReportEnabled" type="String" value="1"/>
44.	Добавлять номер смены к имени файла отчета.
<!-- Add day number to Z report filename - ZReport_0001.xml -->
<prop name="ZReportDayNumber" type="String" value="1"/>
45.	Имя XML файла Z отчета
<!-- XML Z-report file name -->
<prop name="XmlZReportFileName" type="String" value="ZReport.xml"/>
46.	Разрешение получения Z- отчета в виде CSV файла. Драйвер сохраняет данные всех денежных и операционных регистров CSV файле.
<!-- create Z-report in CSV format -->
<prop name="CsvZReportEnabled" type="String" value="0"/>
47.	Имя CSV файла Z отчета
<!-- CSV Z-report file name -->
<prop name="CsvZReportFileName" type="String" value="ZReport.csv"/>
48.	Разрешение обработки ESC команд в тексте. Это нужно для поддержки некоторых устаревших приложений
<!-- ESC commands enabled -->
<prop name="escCommandsEnabled" type="String" value="1"/>
49.	Режим записи таблиц, 0 - автоматически, 1 - запрещена запись в таблицы
<!-- Table mode, 0 - auto, 1 - disabled -->
<prop name="tableMode" type="String" value="0"/>
50.	Режим логотипа перед заголовком чека, 0 - промотка чека на величину клише, 1 - печатать логотип в 2 этапа
<!-- Logo mode, 0 - feed paper, 1 - split image -->
<prop name="logoMode" type="String" value="1"/>
51.	Режим поиска ФР. 0 - нет поиска, 1 - поиск ФР при ошибках
<!-- SearchMode, 0 - none, 1 - search on error -->
<prop name="searchMode" type="String" value="1"/>
52.	Задержка после отрезки чека в миллисекундах
<!-- Paper cut delay in milliseconds -->
<prop name="cutPaperDelay" type="String" value="0"/>
53.	Тип чека продажи, 0 - обычный, 1 - текстовый чек
<!-- Sales receipt type, 0 - normal, 1 - GLOBUS -->
<prop name="salesReceiptType" type="String" value="1"/>
54.	Длина поля "цена" для формата чека
<!-- Amount field length -->
<prop name="RFAmountLength" type="String" value="8"/>
55.	Длина поля "количество" для формата чека
<!-- Quantity field length -->
<prop name="RFQuantityLength" type="String" value="10"/>
56.	Порт для системы мониторинга. Система мониторинга нужна для удаленного запроса состояния ФР и ЭКЛЗ
<!-- Monitoring server port -->
<prop name="MonitoringPort" type="String" value="50000"/>
57.	Разрешение работы системы мониторинга
<!-- Monitoring enabled -->
<prop name="MonitoringEnabled" type="String" value="0"/>
58.	Разрешение сохранения чека в текстовом виде
<!-- Receipt report enabled -->
<prop name="receiptReportEnabled" type="String" value="1"/>
59.	Название файла чека
<!-- Receipt report file name -->
<prop name="receiptReportFileName" type="String" value="ZCheckReport.xml"/>
60.	Тип открытия чека, 0 - открыть чек при печати позиции, 1 - открыть чека в методе beginFiscalReceipt
<!-- openReceiptType, 0 - open receipt on item print, 1 - open receipt in beginFiscalReceipt -->
<prop name="openReceiptType" type="String" value="1"/>
61.	Режим заголовка, 0 - заголовок в драйвере, 1 - в принтере
<!-- headerMode, 0 - header in driver, 1 - header in fiscalprinter -->
<prop name="headerMode" type="String" value="0"/>
62.	Позиция логотипа в заголовке чека
<!-- headerImagePosition,  -->
<!-- SMFPTR_LOGO_AFTER_HEADER     = 0  -->
<!-- SMFPTR_LOGO_BEFORE_TRAILER   = 1  -->
<!-- SMFPTR_LOGO_AFTER_TRAILER    = 2  -->
<!-- SMFPTR_LOGO_AFTER_ADDTRAILER = 3  -->
<!-- SMFPTR_LOGO_BEFORE_HEADER    = 4  -->
<prop name="headerImagePosition" type="String" value="0"/>
63.	Позиция логотипа в рекламном тексте
<!-- trailerImagePosition,  -->
<!-- SMFPTR_LOGO_AFTER_HEADER     = 0  -->
<!-- SMFPTR_LOGO_BEFORE_TRAILER   = 1  -->
<!-- SMFPTR_LOGO_AFTER_TRAILER    = 2  -->
<!-- SMFPTR_LOGO_AFTER_ADDTRAILER = 3  -->
<!-- SMFPTR_LOGO_BEFORE_HEADER    = 4  -->
<prop name="trailerImagePosition" type="String" value="1"/>
64.	Центрировать строки заголовка чека
<!-- Center header and trailer text automatically  -->
<prop name="centerHeader" type="String" value="1"/>
65.	Разрешение лога
<!-- Log file enabled  -->
<prop name="logEnabled" type="String" value="1"/>
66.	Передавать ENQ перед каждой командой
<!-- Send ENQ before every command or not  -->
<prop name="sendENQ" type="String" value="0"/>
67.	Печатать буквы налоговых групп на текстовом чеке
<!-- Enable tax letters for GLOBUS receipt  -->
<prop name="taxLettersEnabled" type="String" value="0"/>
68.	Префикс штрихкода для печати штрихкода как текста
<!-- Barcode prefix -->
<prop name="barcodePrefix" type="String" value="#*~*#"/>
69.	Тип штрихкода: 
UPCA=0, UPCE=1, EAN13=2, EAN8=3, CODE39=4, ITF=5, CODABAR=6, CODE93=7, CODE128=8, PDF417=10, GS1_OMNI=11, GS1_TRUNC=12, GS1_LIMIT=13, GS1_EXP=14, GS1_STK=15, GS1_STK_OMNI=16, GS1_EXP_STK=17, AZTEC=18, DATA_MATRIX=19, MAXICODE=20, QR_CODE=21, RSS_14=22, RSS_EXPANDED=23, UPC_EAN_EXTENSION=24  
<!-- Barcode type -->
<!-- UPCA=0, UPCE=1, EAN13=2, EAN8=3, CODE39=4, ITF=5, CODABAR=6, CODE93=7, CODE128=8, PDF417=10, GS1_OMNI=11,  -->
<!-- GS1_TRUNC=12, GS1_LIMIT=13, GS1_EXP=14, GS1_STK=15, GS1_STK_OMNI=16, GS1_EXP_STK=17, AZTEC=18, DATA_MATRIX=19,  -->
<!-- MAXICODE=20, QR_CODE=21, RSS_14=22, RSS_EXPANDED=23, UPC_EAN_EXTENSION=24  -->
<prop name="barcodeType" type="String" value="21"/>
70.	Ширина штриха в точках. Обычно 2-3 для одномерного штрихкода, 3-4 для двумерного
<!-- Barcode bar/module width -->
<prop name="barcodeBarWidth" type="String" value="4"/>
71.	Высота штрихкода. Имеет значение для одномерных штрихкодов.
<!-- Barcode height -->
<prop name="barcodeHeight" type="String" value="100"/>
72.	Позиция текста относительно штрихкода: 0 - не печатать, 1 - сверху, 2 - снизу, 3 - сверху и снизу
<!-- Barcode text position NOTPRINTED=0, ABOVE=1, BELOW=2, BOTH=3 -->
<prop name="barcodeTextPosition" type="String" value="2"/>
73.	Штрифт текста штрихкода, 1..7
<!-- Barcode text font, 1..7, default 1 -->
<prop name="barcodeTextFont" type="String" value="1"/>
74.	Соотношение ширины и высоты для штрихкода
<!-- Barcode aspect ratio -->
<prop name="barcodeAspectRatio" type="String" value="3"/>
75.	Совместимость со старыми версиями драйвера: 0 - нет, 1 - полная
<!-- Compatibility level, 0 - NONE, 1 - FULL -->
<prop name="compatibilityLevel" type="String" value="0"/>
76.	Задержка после печати штрихкода в миллисекундах
<!-- Delay after barcode printed -->
<prop name="barcodeDelay" type="String" value="0"/>
