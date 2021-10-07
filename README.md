
# logi2 && course work
# in plan

Служба логировагия (журналирования):
- [ ] запись структурированных  сообщений в долгосрочное хранилище (90дней).
- [ ] автоматическая репликация(зеркалирование) без обслуживания.
- [ ] стойкость хранилища к повреждениям/изменениям.
- [ ] удаление сообщений по лимиту времени(ротация)
- [ ] многопользовательский поиск по фразам и параметрам сообщения в том числе нечеткий поиск
- [x] веб интерфейс
- [x] метрики основных процессов(кол-ва сообщений, скорость поиска, входящий поток)

# in future
В дальнейшем:
- [ ] функция забывания. Возможность иметь представление по кол-вам сообщений по срезу времени за долгий период без детализаций

# TODO
- [x] storage xml coded
- [x] read in storage
- [x] write in storage
- [x] UI from storage
- [x] one file bin
- [x] serch on web page (O(n) algorithm complexity) :tada:
- [ ] serch on web page one button :tada:
- [x] add indexing on bleve :tada:
- [x] search on web form use bleve Ulid (O(1) algorithm complexity)
- [ ] create bleve storage for other file 




