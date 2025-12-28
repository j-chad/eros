<script lang="ts">
    import { RefreshCw, Trash2 } from 'lucide-svelte';
    import {api} from "$lib/api";
    import Card from "$lib/components/Card.svelte";
    import Button from "$lib/components/Button.svelte";
    import CopyableCode from "$lib/components/CopyableCode.svelte";
    import Date from "$lib/components/Date.svelte";
    import Table from "$lib/components/Table.svelte";

    let { data } = $props();

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
            await api.devices.delete(deviceId);
            devices = devices.filter(d => d.id !== deviceId);
        }
    }
</script>

<div class="container">
    <h2>Registration</h2>

    <!-- Registration Code Section -->
    <Card title="Registration Code">
        {#snippet actions()}
            <Button variant="primary" onclick={handleRefreshCode}>
                {#snippet icon()}<RefreshCw size={16} />{/snippet}
                Refresh
            </Button>
            {#if registrationCode}
                <Button variant="danger" onclick={handleDeleteCode}>
                    {#snippet icon()}<Trash2 size={16} />{/snippet}
                    Delete
                </Button>
            {/if}
        {/snippet}

        {#if registrationCode}
            <div class="form-group">
                <label>Code</label>
                <CopyableCode code={registrationCode.code} />
            </div>

            <div class="grid-2">
                <div>
                    <label>Created</label>
                    <Date datetime={registrationCode.created_at} inline />
                </div>
                <div>
                    <label>Expires</label>
                    <Date expiry datetime={registrationCode.expires_at} inline />
                </div>
            </div>
        {:else}
            <div class="empty">No active registration code. Click "Refresh" to generate one.</div>
        {/if}
    </Card>

    <!-- Devices Section -->
    <Card title="Registered Devices">
        {#if devices.length > 0}
            <Table headers={['Device Info', 'Registered', 'Last Seen', 'Expires', 'Actions']}>
                {#each devices as device}
                    <tr>
                        <td><div class="device-info">{device.device_info}</div></td>
                        <td class="nowrap">
                            <Date datetime={device.registered_at} />
                        </td>
                        <td class="nowrap">
                            <Date datetime={device.last_seen_at} />
                        </td>
                        <td class="nowrap">
                            <Date datetime={device.expires_at} expiry />
                        </td>
                        <td>
                            <Button
                                    variant="danger"
                                    size="sm"
                                    onclick={() => handleDeleteDevice(device.id)}
                            >
                                {#snippet icon()}<Trash2 size={14} />{/snippet}
                                Remove
                            </Button>
                        </td>
                    </tr>
                {/each}
            </Table>
        {:else}
            <div class="empty">No devices registered yet.</div>
        {/if}
    </Card>
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

    .empty {
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

    .grid-2 {
        display: grid;
        grid-template-columns: 1fr 1fr;
        gap: 1.5rem;
    }

    .info-text {
        font-size: 0.875rem;
        color: #1f2937;
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

    @media (max-width: 768px) {
        .grid-2 {
            grid-template-columns: 1fr;
        }
    }
</style>