import { db, promisifyRequest } from '../db';

export enum KVKey {
	AuthSession = 'auth:session',
}

interface KVValueMap {
	[KVKey.AuthSession]: {
		token: string;
		expiresAt: number;
	};
}

export type KVSchema = {
	[K in KVKey]: KVValueMap[K];
};

export type KVValue<K extends KVKey> = KVSchema[K];

export class KVStore {
	// Set a value with type checking
	static async set<K extends KVKey>(key: K, value: KVValue<K>): Promise<void> {
		const store = await db.getStore('kv', 'readwrite');
		const data = {
			key,
			value,
			timestamp: Date.now()
		};
		await promisifyRequest(store.put(data));
	}

	// Get a value with type inference
	static async get<K extends KVKey>(key: K): Promise<KVValue<K> | null> {
		const store = await db.getStore('kv', 'readonly');
		const result = await promisifyRequest(store.get(key));
		return result ? result as KVValue<K> : null;
	}

	// Delete a value
	static async delete<K extends KVKey>(key: K): Promise<void> {
		const store = await db.getStore('kv', 'readwrite');
		await promisifyRequest(store.delete(key));
	}

	// Check if key exists
	static async has<K extends KVKey>(key: K): Promise<boolean> {
		const store = await db.getStore('kv', 'readonly');
		const result = await promisifyRequest(store.get(key));
		return result !== undefined;
	}

	// Get all keys (typed)
	static async keys(): Promise<KVKey[]> {
		const store = await db.getStore('kv', 'readonly');
		const result = await promisifyRequest(store.getAllKeys());
		return result as KVKey[];
	}

	// Clear all KV data
	static async clear(): Promise<void> {
		const store = await db.getStore('kv', 'readwrite');
		await promisifyRequest(store.clear());
	}

	// Get multiple values at once (typed)
	static async getMany<K extends KVKey>(
		keys: K[]
	): Promise<Map<K, KVValue<K>>> {
		const store = await db.getStore('kv', 'readonly');
		const results = new Map<K, KVValue<K>>();

		await Promise.all(
			keys.map(async (key) => {
				const result = await promisifyRequest(store.get(key));
				if (result) {
					results.set(key, result as KVValue<K>);
				}
			})
		);

		return results;
	}

	// Set multiple values at once (typed)
	static async setMany<K extends KVKey>(
		entries: { [P in K]?: KVValue<P> }
	): Promise<void> {
		const store = await db.getStore('kv', 'readwrite');
		const timestamp = Date.now();

		await Promise.all(
			Object.entries(entries).map(([key, value]) =>
				promisifyRequest(store.put({ key, value, timestamp }))
			)
		);
	}

	// Get with default value
	static async getOrDefault<K extends KVKey>(
		key: K,
		defaultValue: KVValue<K>
	): Promise<KVValue<K>> {
		const value = await this.get(key);
		return value ?? defaultValue;
	}
}
