<!-- src/routes/+page.svelte -->
<script lang="ts">
    import { Eye, EyeOff, Copy, RefreshCw, Trash2, Check } from 'lucide-svelte';
    import type {RegistrationToken} from "$lib/types";
    import {api} from "$lib/api";

    let { data } = $props();

    let codeVisible = $state(false);
    let copied = $state(false);

    let registrationCode = $state(data.registration);
    let devices = $state(data.devices);

    function formatDate(dateString: string) {
        return new Date(dateString).toLocaleString('en-US', {
            year: 'numeric',
            month: 'short',
            day: 'numeric',
            hour: '2-digit',
            minute: '2-digit'
        });
    }

    function getTimeRemaining(expiresAt: string) {
        const now = new Date();
        const expiry = new Date(expiresAt);
        const diff = expiry.getTime() - now.getTime();

        if (diff < 0) return 'Expired';

        const days = Math.floor(diff / (1000 * 60 * 60 * 24));
        const hours = Math.floor((diff % (1000 * 60 * 60 * 24)) / (1000 * 60 * 60));
        const minutes = Math.floor((diff % (1000 * 60 * 60)) / (1000 * 60));

        if (days > 0) return `${days}d ${hours}h`;
        if (hours > 0) return `${hours}h ${minutes}m`;
        return `${minutes}m`;
    }

    async function handleCopyCode() {
        if (!registrationCode) return;

        await navigator.clipboard.writeText(registrationCode.code);
        copied = true;
        setTimeout(() => copied = false, 2000);
    }

    async function handleRefreshCode() {
        registrationCode = await api.registration.refresh();
    }

    async function handleDeleteCode() {
        if (confirm('Are you sure you want to delete this registration code? New devices will not be able to register.')) {
            await api.registration.deleteAll();
            registrationCode = null;
        }
    }

    async function handleDeleteDevice(deviceId: string) {
        if (confirm('Are you sure you want to remove this device? It will need to re-register.')) {
            // Replace with actual API call
            console.log('Deleting device:', deviceId);
            devices = devices.filter(d => d.id !== deviceId);
        }
    }
</script>

<div class="container">
    <h2>Registration</h2>

    <!-- Registration Code Section -->
    <div class="card">
        <div class="card-header">
            <h3>Registration Code</h3>
            <div class="actions">
                <button class="btn btn-primary" onclick={handleRefreshCode}>
                    <RefreshCw size={16} />
                    Refresh
                </button>
                {#if registrationCode}
                    <button class="btn btn-danger" onclick={handleDeleteCode}>
                        <Trash2 size={16} />
                        Delete
                    </button>
                {/if}
            </div>
        </div>

        {#if registrationCode}
            <div class="card-body">
                <div class="form-group">
                    <label>Code</label>
                    <div class="code-input-group">
                        <div class="code-input">
                            {codeVisible ? registrationCode.code : '••••••••••••••••••••••'}
                        </div>
                        <button
                                class="icon-btn"
                                onclick={() => codeVisible = !codeVisible}
                                title={codeVisible ? 'Hide code' : 'Show code'}
                        >
                            {#if codeVisible}
                                <EyeOff size={18} />
                            {:else}
                                <Eye size={18} />
                            {/if}
                        </button>
                        <button
                                class="icon-btn"
                                onclick={handleCopyCode}
                                title="Copy code"
                        >
                            {#if copied}
                                <Check size={18} class="success" />
                            {:else}
                                <Copy size={18} />
                            {/if}
                        </button>
                    </div>
                </div>

                <div class="grid-2">
                    <div>
                        <label>Created</label>
                        <div class="info-text">{formatDate(registrationCode.created_at)}</div>
                    </div>
                    <div>
                        <label>Expires</label>
                        <div class="info-text">
                            {formatDate(registrationCode.expires_at)}
                            <span class="time-remaining">
                                ({getTimeRemaining(registrationCode.expires_at)})
                            </span>
                        </div>
                    </div>
                </div>
            </div>
        {:else}
            <div class="card-body empty">
                No active registration code. Click "Refresh" to generate one.
            </div>
        {/if}
    </div>

    <!-- Devices Section -->
    <div class="card">
        <div class="card-header">
            <h3>Registered Devices</h3>
        </div>

        {#if devices.length > 0}
            <div class="table-container">
                <table>
                    <thead>
                    <tr>
                        <th>Device Info</th>
                        <th>Registered</th>
                        <th>Last Seen</th>
                        <th>Expires</th>
                        <th>Actions</th>
                    </tr>
                    </thead>
                    <tbody>
                    {#each devices as device}
                        <tr>
                            <td>
                                <div class="device-info">{device.deviceInfo}</div>
                            </td>
                            <td class="nowrap">{formatDate(device.registeredAt)}</td>
                            <td class="nowrap">{formatDate(device.lastSeen)}</td>
                            <td class="nowrap">
                                {formatDate(device.expiresAt)}
                                <div class="time-remaining-small">
                                    {getTimeRemaining(device.expiresAt)}
                                </div>
                            </td>
                            <td>
                                <button
                                        class="btn btn-danger btn-sm"
                                        onclick={() => handleDeleteDevice(device.id)}
                                >
                                    <Trash2 size={14} />
                                    Remove
                                </button>
                            </td>
                        </tr>
                    {/each}
                    </tbody>
                </table>
            </div>
        {:else}
            <div class="card-body empty">
                No devices registered yet.
            </div>
        {/if}
    </div>
</div>

<style>
    .container {
        max-width: 100%;
    }

    h2 {
        font-size: 2rem;
        font-weight: bold;
        color: #1f2937;
        margin-bottom: 1.5rem;
    }

    .card {
        background: white;
        border-radius: 8px;
        box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
        margin-bottom: 2rem;
        overflow: hidden;
    }

    .card-header {
        padding: 1rem 1.5rem;
        background: #f9fafb;
        border-bottom: 1px solid #e5e7eb;
        display: flex;
        justify-content: space-between;
        align-items: center;
    }

    .card-header h3 {
        font-size: 1.125rem;
        font-weight: 600;
        color: #1f2937;
        margin: 0;
    }

    .actions {
        display: flex;
        gap: 0.5rem;
    }

    .card-body {
        padding: 1.5rem;
    }

    .card-body.empty {
        text-align: center;
        color: #6b7280;
    }

    .form-group {
        margin-bottom: 1.5rem;
    }

    label {
        display: block;
        font-size: 0.875rem;
        font-weight: 500;
        color: #4b5563;
        margin-bottom: 0.5rem;
    }

    .code-input-group {
        display: flex;
        gap: 0.5rem;
    }

    .code-input {
        flex: 1;
        padding: 0.75rem 1rem;
        background: #f9fafb;
        border: 1px solid #e5e7eb;
        border-radius: 4px;
        font-family: 'Courier New', monospace;
        font-size: 0.875rem;
    }

    .icon-btn {
        padding: 0.75rem;
        background: #f3f4f6;
        border: 1px solid #e5e7eb;
        border-radius: 4px;
        cursor: pointer;
        transition: background 0.2s;
        display: flex;
        align-items: center;
        justify-content: center;
    }

    .icon-btn:hover {
        background: #e5e7eb;
    }

    .icon-btn :global(.success) {
        color: #10b981;
    }

    .grid-2 {
        display: grid;
        grid-template-columns: 1fr 1fr;
        gap: 1.5rem;
    }

    .info-text {
        font-size: 0.875rem;
        color: #1f2937;
    }

    .time-remaining {
        margin-left: 0.5rem;
        color: #f97316;
        font-weight: 500;
    }

    .time-remaining-small {
        font-size: 0.75rem;
        color: #f97316;
        font-weight: 500;
        margin-top: 0.25rem;
    }

    .table-container {
        overflow-x: auto;
    }

    table {
        width: 100%;
        border-collapse: collapse;
    }

    thead {
        background: #f9fafb;
        border-bottom: 1px solid #e5e7eb;
    }

    th {
        padding: 0.75rem 1.5rem;
        text-align: left;
        font-size: 0.75rem;
        font-weight: 500;
        color: #4b5563;
        text-transform: uppercase;
        letter-spacing: 0.05em;
    }

    tbody tr {
        border-bottom: 1px solid #e5e7eb;
        transition: background 0.2s;
    }

    tbody tr:hover {
        background: #f9fafb;
    }

    td {
        padding: 1rem 1.5rem;
        font-size: 0.875rem;
        color: #4b5563;
    }

    td.nowrap {
        white-space: nowrap;
    }

    .device-info {
        font-family: 'Courier New', monospace;
        word-break: break-all;
        max-width: 30rem;
        color: #1f2937;
    }

    .btn {
        display: inline-flex;
        align-items: center;
        gap: 0.5rem;
        padding: 0.5rem 1rem;
        border: none;
        border-radius: 4px;
        cursor: pointer;
        font-size: 0.875rem;
        font-weight: 500;
        transition: background 0.2s;
    }

    .btn-primary {
        background: #3b82f6;
        color: white;
    }

    .btn-primary:hover {
        background: #2563eb;
    }

    .btn-danger {
        background: #ef4444;
        color: white;
    }

    .btn-danger:hover {
        background: #dc2626;
    }

    .btn-sm {
        padding: 0.375rem 0.75rem;
        font-size: 0.875rem;
    }

    @media (max-width: 768px) {
        .grid-2 {
            grid-template-columns: 1fr;
        }

        .card-header {
            flex-direction: column;
            align-items: flex-start;
            gap: 1rem;
        }

        .actions {
            width: 100%;
            flex-direction: column;
        }

        .actions button {
            width: 100%;
        }
    }
</style>