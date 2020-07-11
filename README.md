# rssirc

`rssirc` is an IRC bot which follows RSS feeds.

## Installation

### Using a pre-compiled binary

Download a binary for your system from the releases page.

### From source

You will need a working Go development environment.

```
% go get -u github.com/n7st/rssirc
```

## Configuration

Example configuration:

```yaml
log_level: info
irc:
    channels:
        - "#firstchannel"
        - "#secondchannel"
    server: irc.snoonet.org
    port: 6697
    ident: rssbottest
    max_reconnect: 3
    modes: +B
    nickname: rssbottest
    nickserv_password: my-nickserv-password
    use_tls: true

    # Additional options:
    # debug: true # Enable IRC debugging
rss:
    -
        feed_url: https://www.nasa.gov/rss/dyn/breaking_news.rss
        # poll_delay is in minutes and must be 1 or greater.
        # The default is hourly.
        poll_delay: 60
        # max_history defines the cached RSS feed item length
        max_history: 3
        # channels defines which IRC channels notifications for new feed items
        # will be sent to
        channels:
            - "#firstchannel"
            - "#secondchannel"
```

### Linux

Configuration is read from `$HOME/.config/netsplit/rssirc/config.yaml`.

## Running

```
% rssirc
```

### Optional custom config filename

```
% rssirc path/to/config.yaml
```

## Limitations

* Only new items from followed RSS feeds will be displayed in the IRC channels.
  This is by design and means that the bot can safely be restarted without
  duplicate notifications for new feed items.

## License

MIT. See [LICENSE.txt](./LICENSE.txt).