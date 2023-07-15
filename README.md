# sushi-osc-midi-bridge


Due to some limitations of the midi processing within the ELK Audio Sushi DAW, this
deamon allows to listen to the midi input and converts some of the messages to OSC messages:

| MIDI          | Condition                   | Sushi Command                              |
|---------------|-----------------------------|--------------------------------------------|
| ControlChange | Value < 64 && Control < 10  | Bypass Processor (ID=Control, Value=true)  |
| ControlChange | Value >= 64 && Control < 10 | Bypass Processor (ID=Control, Value=false) |
