<script lang="ts">
	import {FileUp, Plus, RefreshCw, Trash2} from 'lucide-svelte';
    import {api} from "$lib/api";
    import Card from "$lib/components/Card.svelte";
    import Button from "$lib/components/Button.svelte";
    import CopyableCode from "$lib/components/CopyableCode.svelte";
    import DateDisplay from "$lib/components/DateDisplay.svelte";
    import Table from "$lib/components/Table.svelte";
    import EditableField from "$lib/components/EditableField.svelte";
    import Header from "$lib/components/Header.svelte";
	import {buildRegistrationPdf} from "$lib/pdf/registration";

    let {data} = $props();

    let registrationCode = $state(data.registration);
    let devices = $state(data.devices);

    async function handleRefreshDevices() {
        devices = await api.devices.list();
    }

    async function handleCreateCode() {
        registrationCode = await api.registration.create();
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

    async function handleUpdateDeviceInfo(id: string, newValue: string) {
        await api.devices.updateDeviceInfo(id, newValue);

        const device = devices.find(d => d.id === id);
        if (device) {
            device.device_info = newValue;
        }
    }

	async function handleExportPDF() {
		const blob = await buildRegistrationPdf({
			code: registrationCode!.code,
			qrUrl: `${window.location.origin}/register?code=${registrationCode!.code}`,
		})

		const url = URL.createObjectURL(blob);
		window.open(url, '_blank');
	}
</script>

<Header title="Registration Management"/>

<!-- Registration Code Section -->
<Card title="Registration Code">
    {#snippet actions()}
        <Button variant="primary" onclick={handleCreateCode}>
            {#snippet icon()}
                <Plus size={16}/>
            {/snippet}
            New Code
        </Button>
        {#if registrationCode}
            <Button variant="danger" onclick={handleDeleteCode}>
                {#snippet icon()}
                    <Trash2 size={16}/>
                {/snippet}
                Delete
            </Button>
			<Button variant="secondary" onclick={handleExportPDF}>
				{#snippet icon()}
					<FileUp size={16}/>
				{/snippet}
				Export PDF
			</Button>
        {/if}
    {/snippet}

    {#if registrationCode}
        <div class="form-group">
            <label>Code</label>
            <CopyableCode code={registrationCode.code}/>
        </div>

        <div class="grid-2">
            <div>
                <label>Created</label>
                <DateDisplay datetime={registrationCode.created_at} inline/>
            </div>
            <div>
                <label>Expires</label>
                <DateDisplay expiry datetime={registrationCode.expires_at} inline/>
            </div>
        </div>
    {:else}
        <div class="empty">No active registration code. Click "New Code" to generate one.</div>
    {/if}
</Card>

<!-- Devices Section -->
<Card title="Registered Devices">
    {#snippet actions()}
        <Button variant="secondary" onclick={handleRefreshDevices}>
            {#snippet icon()}
                <RefreshCw size={16}/>
            {/snippet}
            Refresh
        </Button>
    {/snippet}
    {#if devices.length > 0}
        <Table headers={['Device Info', 'Registered', 'Last Seen', 'Expires', 'Actions']}>
            {#each devices as device}
                <tr>
                    <td>
                        <div class="device-info">
                            <EditableField bind:value={device.device_info}
                                           onSave={(newValue) => handleUpdateDeviceInfo(device.id, newValue)}
                                           multiline={true}
                            />
                        </div>
                    </td>
                    <td class="nowrap">
                        <DateDisplay datetime={device.registered_at}/>
                    </td>
                    <td class="nowrap">
                        <DateDisplay datetime={device.last_seen_at}/>
                    </td>
                    <td class="nowrap">
                        <DateDisplay datetime={device.expires_at} expiry/>
                    </td>
                    <td>
                        <Button
                                variant="danger"
                                size="sm"
                                onclick={() => handleDeleteDevice(device.id)}
                        >
                            {#snippet icon()}
                                <Trash2 size={14}/>
                            {/snippet}
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

<style>
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
