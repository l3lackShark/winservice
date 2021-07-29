# Go Агент, читающий процессы.

* Cgo требуется только из-за базы данных. Всё чтение выполняется на syscall'ах (пакет memory)
* Запускается через [этот](https://github.com/l3lackShark/winservicerunner) раннер 
* Используется нативный NT путь вместо WIN32, ибо не все приложения могут быть совместимы с win32(Пример: \\wsl$) Если требуется, путь можно сконвертировать https://stackoverflow.com/a/4446189 
* При изменении в процессах в базу отправляется JSON объект:

```json
{
    "new": []Process,
    "closed": []Process
}
```

Время исполнения одного цикла ~~300ms на одном потоке с частотой 3GHz и ~~800ms на частоте 1.3GHz. (Ryzen 4600H)   

# Пример JSON Объекта Process

```json
{
	"name": "taskhostw.exe",
	"pid": 7128,
	"mainModulePath": "\\Device\\HarddiskVolume3\\Windows\\System32\\taskhostw.exe",
	"openTime": "2021-07-23T08:48:32Z",
	"sessionID": 1,
	"sessionUserName": "BlackShark",
	"sessionUserSID": "S-1-5-21-916415770-2301682324-591716892-1001",
	"sessionLoginTime": "2021-07-23T08:48:32Z"
}
```

База данных не имеет никакого шифрования и удаляется после каждого перезапуска сервиса. Пример заполненной таблицы: 
<img  src="https://cdn.discordapp.com/attachments/562954897163812865/870270378151378944/unknown.png"  width="250">


# Сборка
* git clone https://github.com/l3lackShark/winservice
* cd winservice
* go build . (потребуется gcc для sqlite) 

* git clone https://github.com/l3lackShark/winservicerunner
* cd winservicerunner
* go build .
* Изменить пути в `winservicerunner.json` на свои
* Через cmd Администратора установить/запустить сервис: `.\winservicerunner.exe -service install`  `.\winservicerunner.exe -service start`
* Логирование сервиса доступно в "Stderr" и "Stdout" файлах (если есть git bash, можно проследить за файлом через `tail -f filename.log`)
