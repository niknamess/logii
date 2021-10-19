
# logi2 && course work
# in plan

Служба логировагия (журналирования):
- [ ] запись структурированных  сообщений в долгосрочное хранилище (90дней).
- [ ] автоматическая репликация(зеркалирование) без обслуживания.
- [ ] стойкость хранилища к повреждениям/изменениям.
- [ ] удаление сообщений по лимиту времени(ротация)
- [x] многопользовательский поиск по фразам и параметрам сообщения в том числе нечеткий поиск
- [x] веб интерфейс
- [x] метрики основных процессов(кол-ва сообщений, скорость поиска, входящий поток)

# in future
- [ ] функция забывания. Возможность иметь представление по кол-вам сообщений по срезу времени за долгий период без детализаций

# TODO
- [x] storage xml coded
- [x] read in storage
- [x] write in storage
- [x] UI from storage
- [x] one file bin
- [x] serch on web page (O(n) algorithm complexity) :tada:
- [x] serch on web page one button :tada:
- [x] add indexing on bleve :tada:
- [x] search on web form use bleve Ulid (O(1) algorithm complexity)
- [x] create bleve storage for other file :fire:
- [x] js how many clicked so many added files and result fix :cookie:
- [x] split search and indexing string (too long indexing when search file)
- [x] add search in dir 
- [x] indexing file when file changes occur
- [x] add func for compare control sum file
- [ ] add disign in web interface (table)
- [x] add disign in web interface (colors for type message) 




# FLAGS 
    to run use 
    go build && ./logi2 -(flag) something
- "f" parse log file and view in terminal decoded strings from coded XOR logfile ( -f /home/nik/projects/Course/logi2/logtest/test/gen_logs1 )
- "d" parse dir and view in terminal decoded strings from coded XOR logfiles ( -d /home/nik/projects/Course/logi2/logtest/test/ )
- "s" search in coded log files and view in terminal strings ( -s /home/nik/projects/Course/logi2/logtest/test/gen_logs1 )
- "z" server
- "w" write logs to normal in file.txt form from coded files   ( -w /home/nik/projects/Course/logi2/logtest/test/gen_logs1 )
- "g" generate coded logs and write in file ( -g blabla )
- "p" web interface enter port (-p 15000)
- "c" run web_interface and generate log (-c 15000)


# RUN
- make dev 





