export type StopWatch = {
	name: string;

	elapsed: number;
	min: number;
	max: number;
	paused: boolean;
	pauses: number;
	nextPause: number;
};
