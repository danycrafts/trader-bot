export namespace domain {

	export class UserSettings {
	    user_id: number;
	    alpaca_api_key: string;
	    alpaca_secret_key: string;
	    theme: string;
	    notifications_email: boolean;
	    notifications_push: boolean;

	    static createFrom(source: any = {}) {
	        return new UserSettings(source);
	    }

	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.user_id = source["user_id"];
	        this.alpaca_api_key = source["alpaca_api_key"];
	        this.alpaca_secret_key = source["alpaca_secret_key"];
	        this.theme = source["theme"];
	        this.notifications_email = source["notifications_email"];
	        this.notifications_push = source["notifications_push"];
	    }
	}

}
