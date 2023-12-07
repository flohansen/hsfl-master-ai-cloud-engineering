/**
 * Validates EAN (European Article Number) codes including both EAN-13 and EAN-8 formats
 * The code checks if the input is a 13-digit or 8-digit number,
 * calculates and verifies the checksum digit.
 */
export function validateEan(input: number): boolean {
    const inputString = input.toString();

    if (inputString.length === 13) {
        return isValidEAN(inputString, true);
    }

    if (inputString.length === 8) {
        return isValidEAN(inputString, false);
    }

    return false;
}

function isValidEAN(input: string, isLongEan: boolean): boolean {
    const eanRegex: RegExp = isLongEan ? /^[0-9]{13}$/ : /^[0-9]{8}$/;
    if (! eanRegex.test(input)) {
        return false;
    }

    const digits: number[] = input.split('').map(Number);
    const checksum: number | undefined = digits.pop();
    const sum: number = digits.reduceRight((acc: number, digit: number, index: number) => {
        return acc + (isLongEan ? (index % 2 === 0 ? digit : digit * 3) : (index % 2 === 1 ? digit : digit * 3));
    }, 0);

    const calculatedChecksum: number = (10 - (sum % 10)) % 10;
    return checksum === calculatedChecksum;
}
