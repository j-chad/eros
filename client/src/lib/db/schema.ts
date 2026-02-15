export const DB_NAME = 'eros';
export const DB_VERSION = 1;

export interface DBSchema {
	kv: {
		key: string;
		value: {
			key: string;
			value: any;
		}
	};
}

interface StoreParameters extends IDBObjectStoreParameters {
	indexes?: { name: string; keyPath: string | string[]; options?: IDBIndexParameters }[];
}

const STORES:Record<keyof DBSchema, StoreParameters> = {
	kv: {
		keyPath: 'key',
	}
}

export function createObjectStores(db: IDBDatabase) {
	for (const [name, options] of Object.entries(STORES)) {
		if (!db.objectStoreNames.contains(name)) {
			const store = db.createObjectStore(name, options);

			if (options.indexes) {
				for (const index of options.indexes) {
					store.createIndex(index.name, index.keyPath, index.options);
				}
			}

			console.log(`Created object store: ${name}`);
		}
	}
}
