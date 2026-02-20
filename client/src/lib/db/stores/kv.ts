import {db, promisifyRequest} from '../db';

export enum KVKey {
	AuthSession = 'auth:session',
}

interface KVValue {
	[KVKey.AuthSession]: string;
}

export type KVSchema = {
	[K in KVKey]: {
		value: KVValue[K]
		timestamp: number;
	};
};

export class KVStore {
	// Set a value with type checking
	static async set<K extends KVKey>(key: K, value: KVValue[K]): Promise<void> {
		const store = await db.getStore<KVKey, KVSchema[K]>('kv', 'readwrite');
		const data: KVSchema[K] = {
			value,
			timestamp: Date.now()
		};
		await promisifyRequest(store.put(data, key));
	}

	// Get a value with type inference
	static async get<K extends KVKey>(key: K): Promise<KVValue[K] | null> {
		const store = await db.getStore<KVKey, KVSchema[K]>('kv', 'readonly');
		const result = await promisifyRequest(store.get(key));
		return result ? result.value : null;
	}

	// Delete a value
	static async delete<K extends KVKey>(key: K): Promise<void> {
		const store = await db.getStore<KVKey, KVSchema[K]>('kv', 'readwrite');
		await promisifyRequest(store.delete(key));
	}

	// Check if key exists
	static async has<K extends KVKey>(key: K): Promise<boolean> {
		const store = await db.getStore('kv', 'readonly');
		const result = await promisifyRequest(store.get(key));
		return result !== undefined;
	}

	// Get all keys (typed)
	static async keys(): Promise<IDBValidKey[]> {
		const store = await db.getStore('kv', 'readonly');
		return await promisifyRequest(store.getAllKeys());
	}
}
