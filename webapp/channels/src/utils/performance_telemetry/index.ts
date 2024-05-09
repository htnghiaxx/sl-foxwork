// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

export function markAndReport(name: string): PerformanceMark {
    return performance.mark(name, {
        detail: {
            report: true,
        },
    });
}

/**
 * Measures the duration between two performance marks, schedules it to be reported to the server, and returns the
 * PerformanceMeasure created by doing this.
 *
 * If either the start or end mark does not exist, undefined will be returned and, if canFail is false, an error
 * will be logged.
 */
export function measureAndReport(measureName: string, startMark: string, endMark: string, canFail = false): PerformanceMeasure | undefined {
    const options: PerformanceMeasureOptions = {
        start: startMark,
        end: endMark,
        detail: {
            report: true,
        },
    };

    try {
        return performance.measure(measureName, options);
    } catch (e) {
        if (!canFail) {
            // eslint-disable-next-line no-console
            console.error('Unable to measure ' + measureName, e);
        }

        return undefined;
    }
}
