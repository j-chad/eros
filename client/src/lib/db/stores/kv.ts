import {db, promisifyRequest} from '../db';
import type { FavourCount } from '$lib/types/favour';

export enum KVKey {
	AuthSession = 'auth:session',
	CalendarTipSeen = 'tutorial:calendar-tip',
	FavourTipSeen = 'tutorial:favour-tip',
	FavourCount = 'favour:count',
	PushVAPIDKey = 'push:vapid-key'
}

interface KVValue {
	[KVKey.AuthSession]: string;
	[KVKey.CalendarTipSeen]: boolean;
	[KVKey.FavourTipSeen]: boolean;
	[KVKey.FavourCount]: FavourCount;
	[KVKey.PushVAPIDKey]: string;
}

export type KVSchema = {
	[K in KVKey]: {
		key: K;
		value: KVValue[K]
		timestamp: number;
	};
};

export class KVStore {
	// Set a value with type checking
	static async set<K extends KVKey>(key: K, value: KVValue[K]): Promise<void> {
		const store = await db.getStore<KVKey, KVSchema[K]>('kv', 'readwrite');
		const data = { key, value, timestamp: Date.now() } as KVSchema[K];
		await promisifyRequest(store.put(data));
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
