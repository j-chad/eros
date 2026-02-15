import { browser } from '$app/environment';

type StoreNames<DB extends IDBDatabase> = string;

function openDB(dbName: string, version: number, onUpgrade: (db: IDBDatabase) => void) {
	return new Promise<IDBDatabase>((resolve, reject) => {
		const req = indexedDB.open(dbName, version);

		req.onupgradeneeded = () => onUpgrade(req.result);
		req.onsuccess = () => resolve(req.result);
		req.onerror = () => reject(req.error);
	});
}

function reqToPromise<T>(req: IDBRequest<T>) {
	return new Promise<T>((resolve, reject) => {
		req.onsuccess = () => resolve(req.result);
		req.onerror = () => reject(req.error);
	});
}

function txDone(tx: IDBTransaction) {
	return new Promise<void>((resolve, reject) => {
		tx.oncomplete = () => resolve();
		tx.onerror = () => reject(tx.error);
		tx.onabort = () => reject(tx.error);
	});
}

let dbPromise: Promise<IDBDatabase> | null = null;

function getDB() {
	if (!browser) throw new Error('IndexedDB is only available in the browser');
	if (!dbPromise) {
		dbPromise = openDB('eros-app', 1, (db) => {
			if (!db.objectStoreNames.contains('kv')) {
				db.createObjectStore('kv');
			}
		});
	}
	return dbPromise;
}

export async function idbGet<T>(key: string): Promise<T | undefined> {
	const db = await getDB();
	const tx = db.transaction('kv', 'readonly');
	const store = tx.objectStore('kv');
	const value = await reqToPromise(store.get(key));
	await txDone(tx);
	return value as T | undefined;
}

export async function idbSet<T>(key: string, value: T): Promise<void> {
	const db = await getDB();
	const tx = db.transaction('kv', 'readwrite');
	const store = tx.objectStore('kv');
	store.put(value, key);
	await txDone(tx);
}

export async function idbDel(key: string): Promise<void> {
	const db = await getDB();
	const tx = db.transaction('kv', 'readwrite');
	const store = tx.objectStore('kv');
	store.delete(key);
	await txDone(tx);
}
