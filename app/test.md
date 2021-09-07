### How send one message to kkt

0) sendMsgBuf.Write(msg)
    1) defer sendMsgBuf.Clear()
1) -> ENQ (service)
2) <- (ACK | NAK) (service)

    2) ACK
        1) <- MSG (application)
        2) -> ACK (service)
        3) goto 5

    3) NAK
        1) -> (sendMsgBuf.GetMsg()) MSG (application)
        2) <- (ACK | NAK) (service)
            1) ACK - OK (-> ACK) (service)
            2) NAK - WRONG MSG CRC (return err) (service)

### How read one message

1) read stx byte
2) read len byte
3) read bytes defined len
4) check crc