// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import React from 'react';

export default function EmailProduct(): JSX.Element {
    return (
        <div style={{gridArea: 'center', height: '100%', width: '100%'}}>
            <iframe
                title='Foxia Mail'
                src={'https://mail.foxia.vn/?iam_sso=1'}
                style={{border: 'none', width: '100%', height: '100%'}}
            />
        </div>
    );
}


