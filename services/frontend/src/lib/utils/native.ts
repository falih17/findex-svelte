import { Capacitor } from '@capacitor/core';
import { Camera, CameraResultType, CameraSource, CameraDirection } from '@capacitor/camera';
import { Geolocation } from '@capacitor/geolocation';
import { Device } from '@capacitor/device';

export const isNativePlatform = () => Capacitor.isNativePlatform();

/**
 * Capture selfie photo from either Native Camera or Web Media API
 */
export async function takeSelfiePhoto(): Promise<{ blob: Blob; base64?: string }> {
	if (isNativePlatform()) {
		// Use Capacitor Native Camera plugin
		const image = await Camera.getPhoto({
			quality: 85,
			allowEditing: false,
			resultType: CameraResultType.Base64,
			source: CameraSource.Camera,
			direction: CameraDirection.Front, // Front-facing camera for face clock-in
			saveToGallery: false
		});

		if (!image.base64String) {
			throw new Error('Gagal mengambil foto dari kamera');
		}

		// Convert base64 to Blob
		const byteCharacters = atob(image.base64String);
		const byteNumbers = new Array(byteCharacters.length);
		for (let i = 0; i < byteCharacters.length; i++) {
			byteNumbers[i] = byteCharacters.charCodeAt(i);
		}
		const byteArray = new Uint8Array(byteNumbers);
		const blob = new Blob([byteArray], { type: `image/${image.format || 'jpeg'}` });

		return { blob, base64: image.base64String };
	} else {
		throw new Error('Browser mode should use live video canvas stream');
	}
}

/**
 * Get device GPS coordinates with High Accuracy
 */
export async function getCurrentCoordinates(): Promise<{ latitude: number; longitude: number; accuracy: number }> {
	if (isNativePlatform()) {
		// Native Capacitor GPS
		const position = await Geolocation.getCurrentPosition({
			enableHighAccuracy: true,
			timeout: 15000,
			maximumAge: 0
		});

		return {
			latitude: position.coords.latitude,
			longitude: position.coords.longitude,
			accuracy: Math.round(position.coords.accuracy)
		};
	} else {
		// Standard Web Geolocation API
		return new Promise((resolve, reject) => {
			if (!('geolocation' in navigator)) {
				reject(new Error('Geolocation tidak didukung pada browser ini'));
				return;
			}
			navigator.geolocation.getCurrentPosition(
				(pos) => {
					resolve({
						latitude: pos.coords.latitude,
						longitude: pos.coords.longitude,
						accuracy: Math.round(pos.coords.accuracy)
					});
				},
				(err) => reject(new Error(err.message)),
				{ enableHighAccuracy: true, timeout: 15000, maximumAge: 0 }
			);
		});
	}
}

/**
 * Get device information for audit trail & anti-fraud
 */
export async function getDeviceInfo() {
	if (isNativePlatform()) {
		const info = await Device.getInfo();
		const id = await Device.getId();
		return {
			platform: info.platform,
			model: info.model,
			osVersion: info.osVersion,
			deviceId: id.identifier,
			isVirtual: info.isVirtual // Detect Android Emulator
		};
	}
	return {
		platform: 'web',
		model: navigator.userAgent,
		osVersion: 'browser',
		deviceId: 'web-session',
		isVirtual: false
	};
}
