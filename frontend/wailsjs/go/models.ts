export namespace api {
	
	export class IssueKeyResult {
	    id: string;
	    name: string;
	    api_key: string;
	
	    static createFrom(source: any = {}) {
	        return new IssueKeyResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.api_key = source["api_key"];
	    }
	}
	export class ListDecisionsResult {
	    items: db.DecisionItem[];
	    total: number;
	
	    static createFrom(source: any = {}) {
	        return new ListDecisionsResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.items = this.convertValues(source["items"], db.DecisionItem);
	        this.total = source["total"];
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
	export class StatsResult {
	    pending: number;
	    pending_by_level: Record<string, number>;
	    avg_resolution_seconds: number;
	
	    static createFrom(source: any = {}) {
	        return new StatsResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.pending = source["pending"];
	        this.pending_by_level = source["pending_by_level"];
	        this.avg_resolution_seconds = source["avg_resolution_seconds"];
	    }
	}

}

export namespace db {
	
	export class AgentKey {
	    id: string;
	    name: string;
	    // Go type: time
	    created_at: any;
	    // Go type: time
	    revoked_at?: any;
	
	    static createFrom(source: any = {}) {
	        return new AgentKey(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.created_at = this.convertValues(source["created_at"], null);
	        this.revoked_at = this.convertValues(source["revoked_at"], null);
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
	export class DecisionItem {
	    id: string;
	    // Go type: time
	    received_at: any;
	    source: string;
	    agent_id: string;
	    subject: string;
	    level: string;
	    reason: string;
	    summary: string;
	    context: Record<string, any>;
	    status: string;
	    resolution_decision?: string;
	    resolution_feedback?: string;
	    // Go type: time
	    resolution_resolved_at?: any;
	    resolution_resolved_by?: string;
	
	    static createFrom(source: any = {}) {
	        return new DecisionItem(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.received_at = this.convertValues(source["received_at"], null);
	        this.source = source["source"];
	        this.agent_id = source["agent_id"];
	        this.subject = source["subject"];
	        this.level = source["level"];
	        this.reason = source["reason"];
	        this.summary = source["summary"];
	        this.context = source["context"];
	        this.status = source["status"];
	        this.resolution_decision = source["resolution_decision"];
	        this.resolution_feedback = source["resolution_feedback"];
	        this.resolution_resolved_at = this.convertValues(source["resolution_resolved_at"], null);
	        this.resolution_resolved_by = source["resolution_resolved_by"];
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

