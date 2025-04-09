export interface Todo {
  id: number;
  content: string;
}

export interface Meta {
  totalCount: number;
}

export interface MusicTrack {
  name: string
  artist: string
  album: string
  duration: number
  volume: number
}

export interface MusicPlayer {
  track: MusicTrack
  position: number
  paused: boolean
  muted: boolean
  volume: number
}