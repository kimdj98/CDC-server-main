<p align="center">
<img src="https://user-images.githubusercontent.com/88659167/229131308-d658434a-cc34-46d0-a3da-4f2cb86272d0.png" width="100px" alt="HearSitter Logo" />
</p>

# hearsitter-server-go

this code is go server for project hearsitter, solution challenge.

## what this is made of?

this code is made with golang and fiber framework.
also grpc is used for connecting [python server](https://github.com/jimmy0006/hearsitter-server-python).

## how it works?

when application send http request to this server with sound file, this server send the file to python server.
the python server can figure out what kind of sound was send, and can figure out if warning is needed.
this server send sound file to python server with grpc.
after knowing what sound it was, this server sends results to client including sound prediction, status of sync, alarm if warning is needed.

this server just routes sound file, this server can load balance.

![슬라이드1](https://user-images.githubusercontent.com/45549879/225945874-250d63cc-198e-4168-982f-ac4ab5d47274.PNG)

## API reference

|using|method|url|body|response|
|---|---|---|---|---|
|check server alive|get|'/'|-|Hello world!|
|check python server alive|get|'/ping'|-|Pong!|
|send sound file to analyze|post|'/file'|sound file(wav)|{"Alarm":boolean, "Label":String, "Tagging_rate":float}|
|send sound byte[] to analyze|post|'/uint'|sound file as byte[]|{"Alarm":boolean, "Label":String, "Tagging_rate":float}|


if Alarm is true, the warning is needed.
Label tell us which the sound sounds like.

## reference

landing page [https://github.com/mumwa/hearsitter-landingpage](https://github.com/mumwa/hearsitter-landingpage)

go fiber [https://gofiber.io/](https://gofiber.io/)

ml server [https://github.com/jimmy0006/hearsitter-server-ml](https://github.com/jimmy0006/hearsitter-server-ml)
