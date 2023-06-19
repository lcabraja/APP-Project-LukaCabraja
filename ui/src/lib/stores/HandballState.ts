import { writable } from 'svelte/store';
import type { HandballState } from '$lib/model/LiveScore';
import { current_component } from 'svelte/internal';

export const PreviousLiveScoreState = writable<HandballState | null>(null);
export const LiveScoreState = writable<HandballState | null>(null);

export function subscribeHandler(
	currentState: HandballState | null,
	previousState: HandballState | null
) {
	if (currentState === null) return;
	if (previousState === null) return;

	try {
		for (let key in currentState) {
			if (Object(currentState)[key] !== Object(previousState)[key]) {
				console.log(
					`The property ${key} has changed from ${Object(previousState)[key]} to ${
						Object(currentState)[key]
					}`
				);
			}
		}
	} catch (err) {
		console.error(`Error comparing previous state to current state: `, err);
	}

	PreviousLiveScoreState.set(currentState);
}

export function translateState(liveScoreData: any): HandballState {
	let newState: HandballState = {
		matchid: liveScoreData.matchid.value.toString(),
		competition_name: liveScoreData.competition_name.value,
		competition_round: liveScoreData.competition_round.value,
		location: liveScoreData.location.value,
		home: {
			api: {
				name: liveScoreData.api_home_name.value,
				short: liveScoreData.api_home_team.value,
				score: liveScoreData.api_home_score.value,
				scoring: {
					success: liveScoreData.api_home_scoring_success.value,
					percent: liveScoreData.api_home_scoring_percent.value,
					total: liveScoreData.api_home_scoring_total.value
				},
				defense: {
					success: liveScoreData.api_home_defense_success.value,
					percent: liveScoreData.api_home_defense_percent.value,
					total: liveScoreData.api_home_defense_total.value
				},
				fast_break_goals: liveScoreData.api_home_fast_break_goals.value,
				seven_m_goals: liveScoreData.api_home_7m_goals.value,
				lineup: JSON.parse(liveScoreData.api_home_lineup.value),
				coach: JSON.parse(liveScoreData.api_home_coach.value)
			},
			local: {
				name: liveScoreData.home_name.value,
				short: liveScoreData.home_team.value,
				score: liveScoreData.home_score.value,
				scoring: {
					success: liveScoreData.home_scoring_success.value,
					percent: liveScoreData.home_scoring_percent.value,
					total: liveScoreData.home_scoring_total.value
				},
				defense: {
					success: liveScoreData.home_defense_success.value,
					percent: liveScoreData.home_defense_percent.value,
					total: liveScoreData.home_defense_total.value
				},
				fast_break_goals: liveScoreData.home_fast_break_goals.value,
				seven_m_goals: liveScoreData.home_7m_goals.value,
				lineup: JSON.parse(liveScoreData.final_home_lineup.value),
				coach: JSON.parse(liveScoreData.home_coach.value)
			}
		},
		away: {
			api: {
				name: liveScoreData.api_away_name.value,
				short: liveScoreData.api_away_team.value,
				score: liveScoreData.api_away_score.value,
				scoring: {
					success: liveScoreData.api_away_scoring_success.value,
					percent: liveScoreData.api_away_scoring_percent.value,
					total: liveScoreData.api_away_scoring_total.value
				},
				defense: {
					success: liveScoreData.api_away_defense_success.value,
					percent: liveScoreData.api_away_defense_percent.value,
					total: liveScoreData.api_away_defense_total.value
				},
				fast_break_goals: liveScoreData.api_away_fast_break_goals.value,
				seven_m_goals: liveScoreData.api_away_7m_goals.value,
				lineup: JSON.parse(liveScoreData.api_away_lineup.value),
				coach: JSON.parse(liveScoreData.api_away_coach.value)
			},
			local: {
				name: liveScoreData.away_name.value,
				short: liveScoreData.away_team.value,
				score: liveScoreData.away_score.value,
				scoring: {
					success: liveScoreData.away_scoring_success.value,
					percent: liveScoreData.away_scoring_percent.value,
					total: liveScoreData.away_scoring_total.value
				},
				defense: {
					success: liveScoreData.away_defense_success.value,
					percent: liveScoreData.away_defense_percent.value,
					total: liveScoreData.away_defense_total.value
				},
				fast_break_goals: liveScoreData.away_fast_break_goals.value,
				seven_m_goals: liveScoreData.away_7m_goals.value,
				lineup: JSON.parse(liveScoreData.final_away_lineup.value),
				coach: JSON.parse(liveScoreData.away_coach.value)
			}
		},
		time: ''
	};

	return newState;
}
