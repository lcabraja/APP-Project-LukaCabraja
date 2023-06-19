export type Coach = {
	firstname: string;
	lastname: string;
	letter: string;
	personId: number;
	playerId: number;
	team: number;
};

export type Player = {
	firstname: string;
	lastname: string;
	number: number;
};

export type Referee = {
	firstname: string;
	lastname: string;
	order: string;
	personId: string;
	playerId: string;
};

export type TeamState = {
	name: string;
	short: string;
	score: string;
	scoring: {
		success: string;
		percent: string;
		total: string;
	};
	defense: {
		success: string;
		total: string;
		percent: string;
	};
	fast_break_goals: string;
	seven_m_goals: string;
	lineup: [Player];
	coach: Coach;
};

export type HandballState = {
	matchid: string;
	competition_name: string;
	competition_round: string;
	location: string;
	home: {
		api: TeamState;
		local: TeamState;
	};
	away: {
		api: TeamState;
		local: TeamState;
	};
	time: string;
};
