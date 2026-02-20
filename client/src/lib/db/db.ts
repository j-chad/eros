import {createObjectStores, DB_NAME, DB_VERSION, type DBSchema} from "./schema";

interface TypedStore<K extends IDBValidKey, V> extends IDBObjectStore {
	add(value: V, key?: K): IDBRequest<IDBValidKey>;
	get(key: K): IDBRequest<V | undefined>;
	put(value: V, key?: K): IDBRequest<IDBValidKey>;
	delete(key: K): IDBRequest<undefined>;
	getAllKeys(query?: K | IDBKeyRange | null, count?: number): IDBRequest<IDBValidKey[]>;

}

type AutoStore<T extends keyof DBSchema> = TypedStore<DBSchema[T]['key'], DBSchema[T]['value']>;

class Database {
	private dbp: Promise<IDBDatabase> | null = null;

	init(): Promise<IDBDatabase> {
		if (this.dbp) return this.dbp;

		this.dbp = new Promise((resolve, reject) => {
			const request = indexedDB.open(DB_NAME, DB_VERSION);

			request.onerror = () => reject(new Error('Failed to open database'));
			request.onsuccess = () => resolve(request.result);
			request.onupgradeneeded = (event) => {
				console.info(`Upgrading database from version ${event.oldVersion} to ${event.newVersion}`);
				const db = (event.target as IDBOpenDBRequest).result;
				createObjectStores(db);
			}
		});

		return this.dbp;
	}

	async getStore<K extends IDBValidKey = never, V = unknown, T extends keyof DBSchema = any>(
		storeName: T,
		mode: IDBTransactionMode = 'readonly'
	): Promise<
		[K] extends [never]
			? AutoStore<T>
			: TypedStore<K, V>
	> {
		const db = await this.init()
		const transaction = db.transaction(storeName, mode);
		return transaction.objectStore(storeName);
	}

	async transaction<T extends keyof DBSchema>(
		storeNames: T | T[],
		mode: IDBTransactionMode = 'readonly'
	): Promise<IDBTransaction> {
		const db = await this.init();
		return db.transaction(storeNames, mode);
	}

	async close() {
		if (this.dbp) {
			const db = await this.dbp;
			db.close();
			this.dbp = null;
		}
	}
}

export const db = new Database();

export function promisifyRequest<T>(request: IDBRequest<T>): Promise<T> {
	return new Promise((resolve, reject) => {
		request.onsuccess = () => resolve(request.result);
		request.onerror = () => reject(request.error);
	});
}
