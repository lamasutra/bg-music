export namespace model {
	
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

}

