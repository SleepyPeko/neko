import EventEmitter from 'eventemitter3'

export interface ReconnecterAbstractEvents {
  connect: () => void
  disconnect: (error?: Error) => void
}

export abstract class ReconnecterAbstract extends EventEmitter<ReconnecterAbstractEvents> {
  constructor() {
    super()

    if (this.constructor == ReconnecterAbstract) {
      throw new Error("Abstract classes can't be instantiated.");
    }
  }

  public async connect() {
    throw new Error("Method 'connect()' must be implemented.");
  }

  public async disconnect() {
    throw new Error("Method 'disconnect()' must be implemented.");
  }
}

export interface ReconnecterEvents {
  open: () => void
  connect: () => void
  disconnect: () => void
  close: (error?: Error) => void
}

export interface ReconnecterConfig {
  max_reconnects: number
  timeout_ms: number
  backoff_ms: number
}

export class Reconnecter extends EventEmitter<ReconnecterEvents> {
  private _conn: ReconnecterAbstract
  private _config: ReconnecterConfig
  private _timeout: number | undefined

  private _open = false
  private _connected = false
  private _total_reconnects = 0
  private _last_connected: Date | undefined

  constructor(conn: ReconnecterAbstract, config?: ReconnecterConfig) {
    super()

    this._conn = conn
    this._config = {
      max_reconnects: 10,
      timeout_ms: 1500,
      backoff_ms: 750,
      ...config,
    }

    this._conn.on('connect', () => {
      if (this._timeout) {
        window.clearTimeout(this._timeout)
        this._timeout = undefined
      }

      this._connected = true

      if (this._open) {
        this._last_connected = new Date()
        this.emit('connect')
      } else {
        this._conn.disconnect()
      }
    })

    this._conn.on('disconnect', () => {
      if (this._timeout) {
        window.clearTimeout(this._timeout)
        this._timeout = undefined
      }

      this._connected = false

      if (this._open) {
        this.emit('disconnect')
        this.reconnect()
      }
    })
  }

  public get isOpen(): boolean {
    return this._open
  }

  public get isConnected(): boolean {
    return this._connected
  }

  public get totalReconnects(): number {
    return this._total_reconnects
  }

  public get lastConnected(): Date | undefined {
    return this._last_connected
  }

  public open(): void {
    if (this._open) {
      throw new Error('connection is already open')
    }

    this._open = true
    this.emit('open')

    this._conn.connect()
    this._timeout = window.setTimeout(() => {
      this._conn.disconnect()
      this.reconnect()
    }, this._config.timeout_ms)
  }

  public close(error?: Error): void {
    if (!this._open) {
      throw new Error('connection is already closed')
    }

    this._open = false
    this._last_connected = undefined
    this.emit('close', error)

    if (this._connected) {
      this._conn.disconnect()
    }
  }

  public reconnect(): void {
    if (this._connected) {
      throw new Error('connection is already connected')
    }

    this._total_reconnects++

    if (this._config.max_reconnects < this._total_reconnects || this._total_reconnects < 0) {
      setTimeout(this._conn.connect, this._config.backoff_ms)
    } else {
      this.close(new Error('reconnection failed'))
    }
  }
}