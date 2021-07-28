# shtrih-m-driver
go tcp driver for shtrih-m

- Start shift with inn
- Close shift
- Status kkt (healthcheck)
- Check [without print bool]

## TODO: 
- [ ] FN Prints check
- [ ] FN Prints return check

- [ ] config for max timeout request for one connection
- [ ] config for retrys for command

#### *FN(Fiscal accumulator, Cashbox, Фискальный накопитель, Онлайн касса)*

# TODO for v0.2:

1: driver for sending end returning messages to/from Cashbox
return error

2: retry policy for sending messages
if error retry

3: async request for client?
usecase:
client add message for driver; if 
message adding once and always shuld be delivery for FN

client need to know successfull message or unsuccessfull!

4: Ui for settings and FN status
> 
> list FN : Status (Open, Close, Open date) : Name IP : ID : 

