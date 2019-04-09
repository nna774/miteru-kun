# miteru-kun

みてるくんはみてる。

```miteru-kun.service
[Unit]
OnFailure=status-email@%n.service

[Service]
Type=oneshot
ExecStart=/usr/local/bin/miteru-kun -src /bunsyo -dst /data/backup/documents -interval 172800
```

```miterukun.timer
[Timer]
OnCalendar=daily
RandomizedDelaySec=1day
Persistent=true

[Install]
WantedBy=timers.target
```

[status-email@%n.service](https://wiki.archlinux.org/index.php/Systemd/Timers#MAILTO)

src から dst の向きにrsync等で同期が走っていることを期待している時、それが長期間止まっていないことを外形的に監視します。
