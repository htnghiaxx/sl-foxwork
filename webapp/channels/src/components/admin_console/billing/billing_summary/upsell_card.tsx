// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import classNames from 'classnames';
import React from 'react';
import {useIntl} from 'react-intl';

import CloudStartTrialButton from 'components/cloud_start_trial/cloud_start_trial_btn';
import WomanUpArrowsAndCloudsSvg from 'components/common/svg_images_components/woman_up_arrows_and_clouds_svg';
import StartTrialCaution from 'components/pricing_modal/start_trial_caution';

import {openExternalPricingLink, FREEMIUM_TO_ENTERPRISE_TRIAL_LENGTH_DAYS} from 'utils/cloud_utils';
import {t} from 'utils/i18n';
import type {Message} from 'utils/i18n';

import './upsell_card.scss';

const enterpriseAdvantages = [
    {
        id: t('upsell_advantages.onelogin_saml'),
        defaultMessage: 'OneLogin/ADFS SAML 2.0',
    },
    {
        id: t('upsell_advantages.openid'),
        defaultMessage: 'OpenID Connect',
    },
    {
        id: t('upsell_advantages.office365'),
        defaultMessage: 'Office365 suite integration',
    },
];

interface Props {
    advantages: Message[];
    title: Message;
    andMore: boolean;
    cta: Message;
    ctaAction?: () => void;
    ctaPrimary?: boolean;
    upsellIsTrial?: boolean;
}

const andMore = {
    id: t('upsell_advantages.more'),
    defaultMessage: 'And more...',
};

export default function UpsellCard(props: Props) {
    const intl = useIntl();

    const ctaClassname = classNames(
        'UpsellCard__cta',
        {
            btn: props.ctaPrimary,
            'btn-primary': props.ctaPrimary,
        },
    );

    let callToAction = (
        <button
            className={ctaClassname}
            onClick={props.ctaAction}
        >
            {intl.formatMessage(
                {
                    id: props.cta.id,
                    defaultMessage: props.cta.defaultMessage,
                },
                props.cta.values,
            )}
        </button>
    );
    if (props.upsellIsTrial) {
        callToAction = (
            <>
                <CloudStartTrialButton
                    message={
                        intl.formatMessage(
                            {
                                id: props.cta.id,
                                defaultMessage: props.cta.defaultMessage,
                            },
                            props.cta.values,
                        )
                    }
                    telemetryId={'start_cloud_trial_billing_subscription'}
                    extraClass={ctaClassname}
                />
                <p className='disclaimer'>
                    <StartTrialCaution/>
                </p>
            </>
        );
    }
    return (
        <div className='UpsellCard'>
            <div className='UpsellCard__illustration'>
                <WomanUpArrowsAndCloudsSvg
                    width={200}
                    height={200}
                />
            </div>
            <div className='UpsellCard__title'>
                {intl.formatMessage(props.title)}
            </div>
            <div className='UpsellCard__advantages'>
                {props.advantages.map((message: Message) => {
                    return (
                        <div
                            className='advantage'
                            key={message.id}
                        >
                            <i className='fa fa-lock'/>{intl.formatMessage(message)}
                        </div>
                    );
                })}
                {props.andMore && <div className='advantage advantage--more'>
                    <i className='fa fa-lock'/>{intl.formatMessage(andMore)}
                </div>
                }
            </div>
            <div>
                {callToAction}
            </div>
        </div>
    );
}

export const tryEnterpriseCard = (
    <UpsellCard
        title={{
            id: t('admin.billing.subscriptions.billing_summary.try_enterprise'),
            defaultMessage: 'Try Enterprise features for free',
        }}
        cta={{
            id: t('admin.billing.subscriptions.billing_summary.try_enterprise.cta'),
            defaultMessage: 'Try free for {trialLength} days',
            values: {
                trialLength: FREEMIUM_TO_ENTERPRISE_TRIAL_LENGTH_DAYS,
            },
        }}
        andMore={true}
        advantages={enterpriseAdvantages}
        upsellIsTrial={true}
    />
);

export const ExploreEnterpriseCard = () => {
    return (
        <UpsellCard
            title={{

                id: t('admin.billing.subscriptions.billing_summary.explore_enterprise'),
                defaultMessage: 'Explore Enterprise features',
            }}
            cta={{
                id: t('admin.billing.subscriptions.billing_summary.explore_enterprise.cta'),
                defaultMessage: 'View all features',
            }}
            ctaAction={openExternalPricingLink}
            andMore={true}
            advantages={enterpriseAdvantages}
        />
    );
};
