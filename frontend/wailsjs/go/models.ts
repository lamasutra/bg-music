export namespace input {
	
	export class InputDevice {
	    name: string;
	    path: string;
	
	    static createFrom(source: any = {}) {
	        return new InputDevice(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.path = source["path"];
	    }
	}

}

export namespace model {
	
	export class Sfx {
	    volume: number;
	    path: string;
	    sample_rate: number;
	    num_channels: number;
	
	    static createFrom(source: any = {}) {
	        return new Sfx(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.volume = source["volume"];
	        this.path = source["path"];
	        this.sample_rate = source["sample_rate"];
	        this.num_channels = source["num_channels"];
	    }
	}
	export class Event {
	    volume: number;
	    sfx: Sfx[];
	    sentence: string[];
	
	    static createFrom(source: any = {}) {
	        return new Event(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.volume = source["volume"];
	        this.sfx = this.convertValues(source["sfx"], Sfx);
	        this.sentence = source["sentence"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class InputDevice {
	    name: string;
	    path: string;
	
	    static createFrom(source: any = {}) {
	        return new InputDevice(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.path = source["path"];
	    }
	}
	export class InputDeviceControls {
	    device: InputDevice;
	    mapping: Record<string, string>;
	
	    static createFrom(source: any = {}) {
	        return new InputDeviceControls(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.device = this.convertValues(source["device"], InputDevice);
	        this.mapping = source["mapping"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class Music {
	    volume: number;
	    path: string;
	    skip: number;
	    endBefore: number;
	
	    static createFrom(source: any = {}) {
	        return new Music(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.volume = source["volume"];
	        this.path = source["path"];
	        this.skip = source["skip"];
	        this.endBefore = source["endBefore"];
	    }
	}
	export class MusicMetadata {
	    artist: string;
	    title: string;
	    album: string;
	    year: number;
	    duration: number;
	
	    static createFrom(source: any = {}) {
	        return new MusicMetadata(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.artist = source["artist"];
	        this.title = source["title"];
	        this.album = source["album"];
	        this.year = source["year"];
	        this.duration = source["duration"];
	    }
	}
	
	export class State {
	    volume?: number;
	    music: Music[];
	    states: string[];
	
	    static createFrom(source: any = {}) {
	        return new State(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.volume = source["volume"];
	        this.music = this.convertValues(source["music"], Music);
	        this.states = source["states"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}

}

export namespace ui {
	
	export class WindowSize {
	    width: number;
	    height: number;
	
	    static createFrom(source: any = {}) {
	        return new WindowSize(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.width = source["width"];
	        this.height = source["height"];
	    }
	}

}

