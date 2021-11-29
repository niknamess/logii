
# Текущая структура сообщения СТКУ

struct Record {

enum class Type {
    FIRST,

    Info = FIRST,
    Debug,
    Warning,
    Error,
    Fatal,

    COUNT
};

QString app_name;
QString app_path;
qint64 app_pid = 0;
QString thread;
QDateTime time;
quint8 type = 0;
QString message;
QString details;
QString address;
QString ulid;

};
