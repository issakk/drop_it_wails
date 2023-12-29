export namespace model {

    export class FileInfo {
        size: number;
        name: string;
        date: string;

        constructor(source: any = {}) {
            if ('string' === typeof source) source = JSON.parse(source);
            this.size = source["size"];
            this.name = source["name"];
            this.date = source["date"];
        }

        static createFrom(source: any = {}) {
            return new FileInfo(source);
        }
    }
	export class TreeNode {
        value: string;
        label: string;
        children?: TreeNode[];

        constructor(source: any = {}) {
            if ('string' === typeof source) source = JSON.parse(source);
            this.value = source["value"];
            this.label = source["label"];
            this.children = this.convertValues(source["children"], TreeNode);
        }

        static createFrom(source: any = {}) {
            return new TreeNode(source);
        }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
            if (!a) {
                return a;
            }
            if (a.slice) {
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

